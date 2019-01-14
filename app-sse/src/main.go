package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

type server struct {
	clt      pubsub.SubscriberClient
	timeSSE  *prometheus.HistogramVec
	messages chan libmetier.MessageSocial
	b        *Broker
}

var (
	subName = flag.String("sub-name", os.Getenv("SUB_NAME"), "the pubsub sunbscription")
	message pubsub.PubsubMessage
)

func main() {

	flag.Parse()

	var s server

	// Sub client
	s.clt = connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	s.timeSSE = promHistogramVec()

	s.messages = make(chan libmetier.MessageSocial)

	// Make a new Broker instance
	s.b = &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	// Consume message on the sub
	log.Println("launch consume thread")
	go s.consumemessage()

	// Start processing events
	log.Println("Start processing events")
	go s.b.start()

	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.
	log.Println("Start get messages function")
	go func() {
		for {
			ms := <-s.messages
			b, err := json.Marshal(ms)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			for i := 0; i != len(s.b.clients); i++ {
				s.b.messages <- fmt.Sprintf(string(b))
			}
		}
	}()

	log.Println("Start end function")
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 1 second to finish processing")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	http.Handle("/", http.HandlerFunc(handler))
	http.Handle("/events/", s.b)
	http.Handle("/metrics", promhttp.Handler())
	log.Println("launch server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

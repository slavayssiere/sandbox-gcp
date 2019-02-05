package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
	sub      string
	ctx      context.Context
}

var (
	topicName = flag.String("topic-name", os.Getenv("TOPIC_NAME"), "the pubsub subscription")
	message   pubsub.PubsubMessage
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	flag.Parse()

	var s server

	s.ctx = context.Background()

	sha256 := sha256.Sum256([]byte(time.Now().Format(time.RFC1123)))

	s.sub = fmt.Sprintf("projects/slavayssiere-sandbox/subscriptions/app-sse-subcription-%x", sha256)

	// Sub client
	s.clt = connexionSubcriber(s.ctx, s.sub, "pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

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
		s.closeSubscription()
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

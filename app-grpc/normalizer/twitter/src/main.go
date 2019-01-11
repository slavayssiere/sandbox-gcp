package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
)

var (
	addr    = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	topic   = flag.String("topic-name", os.Getenv("TOPIC_NAME"), "The topic listen.")
	subname = flag.String("sub-name", os.Getenv("SUB_NAME"), "the subscription write")
)

type server struct {
	pub         pubsub.PublisherClient
	sub         pubsub.SubscriberClient
	tweetStream chan twitter.Tweet
	msgStream chan libmetier.MessageSocial
	timeProm        *prometheus.HistogramVec
}

func (s server) convert() {
	for {
		tweet <- s.tweetStream 
		var u libmetier.MessageSocial
		u.Data = tweet.Text
		u.User = tweet.User.Email
		u.Source = "twitter"
		s.msgStream <- u
	}
}

func main() {

	var s server

	rand.Seed(time.Now().UnixNano())
	s.pub = connexionPublisher("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	s.sub = connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	s.tweetStream = make(chan twitter.Tweet)

	s.timeProm = getPromTime()

	log.Println("launch converter thread")
	go s.convert()

	println("launch consume thread")
	go s.consumemessage()

	println("launch send thread")
	go s.sendMessage()

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

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

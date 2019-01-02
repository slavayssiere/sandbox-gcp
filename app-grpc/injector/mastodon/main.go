package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mastodon "github.com/mattn/go-mastodon"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"context"
)

var (
	addr = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
)

func connexionPublisher(address string, filename string, scope ...string) pubsub.PublisherClient {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	fmt.Printf("Secret in %s\n", filename)
	perRPC, _ := oauth.NewServiceAccountFromFile(filename, "https://www.googleapis.com/auth/pubsub")
	conn, _ := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)

	return pubsub.NewPublisherClient(conn)
}

func publishmessage(maStatus *mastodon.Status, client pubsub.PublisherClient, publishTime chan int64) {
	var message pubsub.PubsubMessage
	var request pubsub.PublishRequest

	start := time.Now()
	ctx := context.Background()

	message.Data = []byte(maStatus.Content)
	message.Attributes = make(map[string]string)
	message.Attributes["source"] = "mastodon"
	message.Attributes["time"] = strconv.FormatInt(start.UnixNano(), 10)

	request.Topic = os.Getenv("TOPIC_NAME")
	request.Messages = append(request.Messages, &message)

	if _, err := client.Publish(ctx, &request); err != nil {
		fmt.Println(err)
		println("error")
	}

	t := time.Now()
	elapsed := t.Sub(start)

	publishTime <- elapsed.Nanoseconds()
}

func main() {
	// Client
	clientPub := connexionPublisher("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	clientMastodon := mastodon.NewClient(&mastodon.Config{
		Server:       os.Getenv("SERVER"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	})

	err := clientMastodon.Authenticate(context.Background(), os.Getenv("LOGIN"), os.Getenv("PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	// Prometheus
	// histogramMean := PromHistogramVec()
	// messagesCounter := PromCounterVec()
	publishTime := make(chan int64)
	// println("Launch mean calculation thread")
	// go func() {
	// 	for {
	// 		elapsed := <-publishTime
	// 		histogramMean.WithLabelValues(os.Getenv("MESSAGE_SIZE"), os.Getenv("TOPIC_NAME")).Observe(float64(elapsed))
	// 	}
	// }()

	timeline, err := clientMastodon.StreamingHashtag(context.Background(), os.Getenv("HASHTAG"), false)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for e := range timeline {
			if _, ok := e.(*mastodon.ErrorEvent); !ok {
				publishmessage(e.(*mastodon.UpdateEvent).Status, clientPub, publishTime)
			}
		}
	}()

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))

	// graceful
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
}

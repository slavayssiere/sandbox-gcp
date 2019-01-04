package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"context"
)

var addr = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")

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

func publishmessage(arraySize int, client pubsub.PublisherClient, publishTime chan int64) {
	start := time.Now()

	var message pubsub.PubsubMessage
	dataMessage := make([]byte, arraySize)
	rand.Read(dataMessage)
	message.Data = dataMessage
	message.Attributes = make(map[string]string)
	message.Attributes["time"] = strconv.FormatInt(start.UnixNano(), 10)

	var request pubsub.PublishRequest
	request.Topic = os.Getenv("TOPIC_NAME")
	request.Messages = append(request.Messages, &message)

	ctx := context.Background()

	if _, err := client.Publish(ctx, &request); err != nil {
		fmt.Println(err)
		println("error")
	}

	t := time.Now()
	elapsed := t.Sub(start)

	publishTime <- elapsed.Nanoseconds()
}

func main() {

	clientPub := connexionPublisher("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	histogramMean := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "mean_in_injector",
		Help: "Time for pubish to pubsub in nanosecond",
	}, []string{"size", "trade"})

	messagesCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_injected",
			Help: "How many messages injected, partitioned by size and trade",
		},
		[]string{"size", "trade"},
	)
	prometheus.Register(histogramMean)
	prometheus.Register(messagesCounter)

	publishTime := make(chan int64)

	println("Launch mean calculation thread")
	go func() {
		for {
			elapsed := <-publishTime
			histogramMean.WithLabelValues(os.Getenv("MESSAGE_SIZE"), os.Getenv("TOPIC_NAME")).Observe(float64(elapsed))
		}
	}()

	msgsize, err := strconv.Atoi(os.Getenv("MESSAGE_SIZE"))
	if err != nil {
		fmt.Println(err)
		return
	}
	freqpers, err := strconv.Atoi(os.Getenv("FREQUENCY_PER_SECOND"))
	if err != nil {
		fmt.Println(err)
		return
	}

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

	println("Launch injection threads")
	for i := 0; i < freqpers; i++ {
		go func() {
			for {
				startTime := time.Now()
				publishmessage(msgsize, clientPub, publishTime)
				messagesCounter.WithLabelValues(strconv.Itoa(msgsize), os.Getenv("TOPIC_NAME")).Add(1)
				elapsedTime := time.Since(startTime)
				time.Sleep((1000 * time.Millisecond) - (elapsedTime / time.Millisecond))
			}
		}()
	}

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

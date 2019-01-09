package main

import (
	"bytes"
	"crypto/x509"
	"encoding/gob"
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
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"context"
)

// Entry test
type Entry struct {
	EntryTime  int64
	ID         string
	Datasource string
	Field1     int64
	Field2     int64
	StartTime  int64
}

func encodeEntry(entry *Entry) *bytes.Buffer {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(entry)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return &b
}

var addr = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")

func connexionPublisher(address string, filename string, scope ...string) pubsub.PublisherClient {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := oauth.NewServiceAccountFromFile(filename, scope...)
	//perRPC := oauth.NewComputeEngine()
	conn, _ := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)

	return pubsub.NewPublisherClient(conn)
}

func connexionSubcriber(address string, filename string, scope ...string) pubsub.SubscriberClient {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := oauth.NewServiceAccountFromFile(filename, scope...)
	conn, _ := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)

	return pubsub.NewSubscriberClient(conn)
}

func publishmessage(mess *pubsub.PubsubMessage, client pubsub.PublisherClient) {

	var request pubsub.PublishRequest
	request.Topic = os.Getenv("TOPIC_NAME")
	request.Messages = append(request.Messages, mess)

	fmt.Println(request)

	ctx := context.Background()

	if _, err := client.Publish(ctx, &request); err != nil {
		fmt.Println(err)
	}
}

func createEntry(starttime int64) *Entry {
	entry := Entry{
		EntryTime:  time.Now().UnixNano(),
		ID:         strconv.Itoa(rand.Int()),
		Datasource: os.Getenv("TOPIC_NAME"),
		Field1:     rand.Int63n(1000),
		Field2:     2,
		StartTime:  starttime,
	}
	return &entry
}

func processEntries(data []byte, entryStream chan *Entry, starttime int64) {
	entryCount := len(data) / 128
	for i := 0; i < entryCount; i++ {
		var entry = createEntry(starttime)
		entryStream <- entry
	}
}

func consumemessage(client pubsub.SubscriberClient, entryStream chan *Entry) {

	var pull pubsub.PullRequest
	pull.Subscription = os.Getenv("SUB_NAME")
	pull.MaxMessages = 5

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "time_in_normalizer",
		Help: "Time for normalizer project in nanosecond",
	}, []string{"normalizer"})

	prometheus.Register(histogram)

	for {
		ctx := context.Background()
		resp, err := client.Pull(ctx, &pull)
		if err != nil {
			fmt.Println(err)
		} else {
			var ackMess pubsub.AcknowledgeRequest
			ackMess.Subscription = pull.Subscription
			for _, messRec := range resp.ReceivedMessages {
				ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
				mess := messRec.GetMessage()
				if starttime, err := strconv.ParseInt(mess.Attributes["time"], 10, 64); err != nil {
					fmt.Println(err)
				} else {
					var elapsedTime float64
					elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano() - starttime)
					histogram.WithLabelValues("time").Observe(elapsedTime)
					processEntries(mess.Data, entryStream, starttime)
				}

			}
			// fmt.Println(ackMess.AckIds)
			client.Acknowledge(ctx, &ackMess)
		}
	}
}

func sendMessage(client pubsub.PublisherClient, entryStream chan *Entry) {

	// defer wg.Done()

	for {
		entry := <-entryStream

		var message pubsub.PubsubMessage
		var encodedEntry = encodeEntry(entry)
		message.Data = encodedEntry.Bytes()
		message.Attributes = make(map[string]string)
		message.Attributes["time"] = strconv.FormatInt(entry.StartTime, 10)
		message.Attributes["normalizer_time"] = strconv.FormatInt(time.Now().UnixNano(), 10)

		publishmessage(&message, client)
	}
}

func main() {

	// rand.Seed(time.Now().UnixNano())
	clientPub := connexionPublisher("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	clientSub := connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	var entryStream chan *Entry
	entryStream = make(chan *Entry)

	println("launch consume thread")
	go consumemessage(clientSub, entryStream)

	println("launch send thread")
	go sendMessage(clientPub, entryStream)

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

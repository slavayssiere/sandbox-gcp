package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func createEntry(starttime int64, data *pubsub.PubsubMessage) *Entry {
	entry := Entry{
		// random time between 2017-01-01 and 2018-01-01 in 100 ms intervals
		EntryTime:  1483228800000 + rand.Int63n(3153600000),
		ID:         strconv.Itoa(rand.Int()),
		Datasource: os.Getenv("TOPIC_NAME"),
		Source:     data.Attributes["source"],
		Text:       data.Data,
		StartTime:  starttime,
	}
	return &entry
}

func processEntries(data *pubsub.PubsubMessage, entryStream chan *Entry, starttime int64) {
	var entry = createEntry(starttime, data)
	entryStream <- entry
}

func consumemessage(client pubsub.SubscriberClient, entryStream chan *Entry) {

	var pull pubsub.PullRequest
	pull.Subscription = os.Getenv("SUB_NAME")
	pull.MaxMessages = 5

	ctx := context.Background()

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "time_in_consumer",
		Help: "Time for consumer project",
	}, []string{"consumer"})

	prometheus.Register(histogram)

	for {
		if resp, err := client.Pull(ctx, &pull); err != nil {
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
					elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano()-starttime) / float64(time.Millisecond)
					histogram.WithLabelValues("complete_time").Observe(elapsedTime)
				}
				if starttime, err := strconv.ParseInt(mess.Attributes["time_normalize"], 10, 64); err != nil {
					fmt.Println(err)
				} else {
					var elapsedTime float64
					elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano()-starttime) / float64(time.Millisecond)
					histogram.WithLabelValues("normalize_time").Observe(elapsedTime)
					processEntries(mess, entryStream, starttime)
				}
			}
			client.Acknowledge(ctx, &ackMess)
		}
	}
}

func connexionSubcriber(address string, filename string, scope ...string) pubsub.SubscriberClient {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := oauth.NewServiceAccountFromFile(filename, "https://www.googleapis.com/auth/pubsub")
	conn, _ := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)

	return pubsub.NewSubscriberClient(conn)
}

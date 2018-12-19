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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"context"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

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

func consumemessage(client pubsub.SubscriberClient) {

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
				}
			}
			client.Acknowledge(ctx, &ackMess)
		}
	}
}

func main() {

	clientSub := connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	go consumemessage(clientSub)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigtable"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"

	"context"
)

var (
	addr           = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	hashtag        = flag.String("hashtag", os.Getenv("HASHTAG"), "Twitter hashtag")
	projectid      = flag.String("project-id", os.Getenv("PROJECT_ID"), "Twitter hashtag")
	instanceid     = flag.String("instance-id", os.Getenv("INSTANCE_ID"), "Twitter hashtag")
	tableid     = flag.String("table-id", os.Getenv("TABLE_ID"), "Twitter hashtag")
	subname        = flag.String("sub-name", os.Getenv("SUB_NAME"), "Twitter hashtag")
)

type server struct {
	sub pubsub.SubscriberClient
	bt bigtable.Client
	messages chan libmetier.MessageSocial
	timeProm *prometheus.HistogramVec
}

func main() {

	flag.Parse()
	var s server

	// Define globals
	ctx := context.Background()

	s.sub = s.connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	s.bt = bigtableClient(ctx)
	s.messages = make(chan libmetier.MessageSocial)
	s.timeProm = promHistogramVec()

	println("launch consume thread")
	go s.consumemessage()

	println("write in bigtable")
	go s.writeMessages(ctx)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

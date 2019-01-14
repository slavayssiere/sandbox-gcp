package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
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
	subname        = flag.String("sub-name", os.Getenv("SUB_NAME"), "Twitter hashtag")
	secretpath     = flag.String("secret-path", os.Getenv("SECRET_PATH"), "Twitter hashtag")
)

type server struct {
	sub pubsub.SubscriberClient
	ds *datastore.Client
	messages chan libmetier.MessageSocial
	timeProm *prometheus.HistogramVec
}

func main() {

	flag.Parse()
	var s server

	// Define globals
	ctx := context.Background()

	log.Println("Get secret from: " + *secretpath)
	s.sub = s.connexionSubcriber("pubsub.googleapis.com:443", *secretpath, "https://www.googleapis.com/auth/pubsub")
	s.ds = datastoreClient(ctx)
	s.messages = make(chan libmetier.MessageSocial)
	s.timeProm = promHistogramVec()

	println("launch consume thread")
	go s.consumemessage()

	println("write in bigtable")
	go s.writeMessages(ctx)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

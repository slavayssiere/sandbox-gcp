package main

import (
	"github.com/go-redis/redis"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"context"
)

// LoggerMiddleware add logger and metrics
func LoggerMiddleware(inner http.HandlerFunc, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		inner.ServeHTTP(w, r)

		if strings.Compare(name,"health") != 0 {
			time := time.Since(start)
			log.Printf(
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				name,
				time,
			)
		}
	})
}

var (
	addr           = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	hashtag        = flag.String("hashtag", os.Getenv("HASHTAG"), "Twitter hashtag")
	projectid      = flag.String("project-id", os.Getenv("PROJECT_ID"), "Twitter hashtag")
	subname        = flag.String("sub-name", os.Getenv("SUB_NAME"), "Twitter hashtag")
	secretpath     = flag.String("secret-path", os.Getenv("SECRET_PATH"), "Twitter hashtag")
	redisaddr      = flag.String("redis-address", os.Getenv("REDIS_HOST")+":6379", "The address to listen on for HTTP requests.")
)

type server struct {
	sub pubsub.SubscriberClient
	ds *datastore.Client
	messages chan libmetier.MessageSocial
	timeProm *prometheus.HistogramVec
	redis *redis.Client
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
	s.redis = redisNew()

	println("launch consume thread")
	go s.consumemessage()

	println("write in bigtable")
	go s.writeMessages(ctx)

	router := mux.NewRouter().StrictSlash(true)

	var handlerStatus http.Handler
	handlerStatus = LoggerMiddleware(libmetier.HandlerStatusFunc, "root")
	router.
		Methods("GET").
		Path("/").
		Name("root").
		Handler(handlerStatus)

	var handlerHealth http.Handler
	handlerHealth = LoggerMiddleware(libmetier.HandlerHealthFunc, "health")
	router.
		Methods("GET").
		Path("/healthz").
		Name("health").
		Handler(handlerHealth)
	
	var handlerUsers http.Handler
	handlerUsers = LoggerMiddleware(handlerUsersFunc, "users_get")
	router.
		Methods("GET").
		Path("/users").
		Name("users_get").
		Handler(handlerUsers)

	router.Methods("GET").Path("/metrics").Name("Metrics").Handler(promhttp.Handler())

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Printf("Start server...")
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))
}

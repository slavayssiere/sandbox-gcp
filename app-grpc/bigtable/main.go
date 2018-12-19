package main

import (
	"bytes"
	"crypto/x509"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"context"
)

// Entry struct test
type Entry struct {
	entryTime  int64
	ID         string
	Datasource string
	field1     int64
	field2     int64
	StartTime  int64
}

var addr = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")

func bigtableClient(ctx context.Context, projectID string, instanceID string) bigtable.Client {
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}

	return *client
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

func decodeEntry(b *bytes.Buffer) *Entry {
	var entry Entry
	d := gob.NewDecoder(b)
	if err := d.Decode(&entry); err != nil {
		panic(err)
	}
	return &entry
}

func intToBytes(input int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, input)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func writeMessage(ctx context.Context, mess *pubsub.PubsubMessage, client bigtable.Client) {

	tableID := os.Getenv("TABLE_ID")

	tbl := client.Open(tableID)

	var entry = decodeEntry(bytes.NewBuffer(mess.Data))

	ds := strings.Split(entry.Datasource, "/")

	// Mutation Way
	mut := bigtable.NewMutation()
	mut.Set("tests", "entry_time", bigtable.Now(), intToBytes(entry.entryTime))
	mut.Set("tests", "id", bigtable.Now(), []byte(entry.ID))
	mut.Set("tests", "datasource", bigtable.Now(), []byte(ds[3]))
	mut.Set("tests", "field1", bigtable.Now(), intToBytes(entry.field1))
	mut.Set("tests", "field2", bigtable.Now(), intToBytes(entry.field2))

	// Read pubsub attribute key to determine BT row key
	var key = entry.Datasource + "_" + strconv.Itoa(int(entry.entryTime)) + "_" + entry.ID

	if err := tbl.Apply(ctx, key, mut); err != nil {
		fmt.Println(err)
	}
}

func consumeMessage(client pubsub.SubscriberClient, messageStream chan *pubsub.PubsubMessage) {

	var pull pubsub.PullRequest
	pull.Subscription = os.Getenv("SUB_NAME")
	pull.MaxMessages = 5

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "time_in_bigtable",
		Help: "Time for bigtable project in nanosecond",
	}, []string{"bigtable"})

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
					elapsedTime = float64(time.Now().UnixNano() - starttime)
					histogram.WithLabelValues("time").Observe(elapsedTime)
				}

				messageStream <- mess
			}
			// fmt.Println(ackMess.AckIds)
			client.Acknowledge(ctx, &ackMess)
		}
	}
}

func sendMessage(ctx context.Context, client bigtable.Client, messageStream chan *pubsub.PubsubMessage) {

	for {

		mess := <-messageStream
		writeMessage(ctx, mess, client)
	}
}

func main() {

	// Define globals
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	instanceID := os.Getenv("INSTANCE_ID")
	secretPath := os.Getenv("SECRET_PATH")

	clientSub := connexionSubcriber("pubsub.googleapis.com:443", secretPath, "https://www.googleapis.com/auth/pubsub")

	clientBT := bigtableClient(ctx, projectID, instanceID)

	var messageStream chan *pubsub.PubsubMessage
	messageStream = make(chan *pubsub.PubsubMessage)

	println("launch consume thread")
	go consumeMessage(clientSub, messageStream)

	println("write in bigtable")
	go sendMessage(ctx, clientBT, messageStream)

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

type Entry struct {
	EntryTime  int64
	ID         string
	Datasource string
	Source     string
	Text       []byte
	StartTime  int64
}

type Broker struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

var message pubsub.PubsubMessage

func encodeEntry(entry *Entry) *bytes.Buffer {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(entry)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return &b
}

// This Broker method starts a new goroutine.  It handles
// the addition & removal of clients, as well as the broadcasting
// of messages out to clients that are currently attached.
//
func (b *Broker) Start() {

	go func() {
		for {

			// look the different event on /events
			select {

			case s := <-b.newClients:

				// There is a new client attached and we
				// want to start sending them messages.
				b.clients[s] = true
				log.Println("Added new client")

			case s := <-b.defunctClients:

				// delete client
				delete(b.clients, s)
				close(s)

				log.Println("Removed client")

			case msg := <-b.messages:

				// send message
				for s := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

// This Broker method handles and HTTP request at the "/events/" URL.
//
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	//
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a new channel, over which the broker can
	// send this client messages.
	messageChan := make(chan string)

	// Add this client to the map of those that should
	// receive updates
	b.newClients <- messageChan

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {

		// Read from our messageChan.
		msg, open := <-messageChan

		if !open {
			// If our messageChan was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		fmt.Fprintf(w, "data: Message: %s\n\n", msg)

		// Flush the response.  This is only possible if
		// the repsonse supports streaming.
		f.Flush()
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}

// handler http (index /)
func handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("WTF dude, error parsing your template.")

	}

	t.Execute(w, os.Getenv("TOPIC_NAME"))
	log.Println("Finished HTTP request at", r.URL.Path)
}

func main() {

	var entryStream chan *Entry

	// Sub client
	clientSub := connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	entryStream = make(chan *Entry)

	// Make a new Broker instance
	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	// Consume message on the sub
	fmt.Println("launch consume thread")
	go consumemessage(clientSub, entryStream)

	// Start processing events
	b.Start()
	http.Handle("/events/", b)

	go func() {
		for {
			entry := <-entryStream

			var encodedEntry = encodeEntry(entry)
			message.Data = encodedEntry.Bytes()
			fmt.Println(message)
		}
	}()

	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.
	go func() {
		for i := 0; ; i++ {
			// send message in the HTTP stream
			b.messages <- fmt.Sprintf(string(message.Data))
			log.Printf("Sent message %d ", i)
			time.Sleep(5e9)
		}
	}()

	http.Handle("/", http.HandlerFunc(handler))
	http.ListenAndServe(":8000", nil)
}

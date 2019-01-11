package main

import (
	"strconv"
	"time"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

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

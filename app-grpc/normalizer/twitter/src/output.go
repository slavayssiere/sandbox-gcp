package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

func (s server) sendMessage() {

	for {
		msg := <-s.msgStream

		var message pubsub.PubsubMessage
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
		ctx := context.Background()
		message.Data = []byte(b)
		message.Attributes = make(map[string]string)
		message.Attributes["time"] = strconv.FormatInt(entry.StartTime, 10)
		message.Attributes["source"] = "twitter"

		s.publishmessage(&message)
	}
}

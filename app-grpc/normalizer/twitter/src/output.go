package main

import (
	"encoding/json"
	"log"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

func (s server) sendMessage() {

	for {
		log.Println("Wait for msgSTream...")
		msg := <-s.msgStream

		var message pubsub.PubsubMessage
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
		message.Data = []byte(b)
		message.Attributes = make(map[string]string)
		message.Attributes["source"] = "twitter"

		s.publishmessage(&message)
	}
}

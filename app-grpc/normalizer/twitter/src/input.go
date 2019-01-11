package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

func (s server) consumemessage() {

	var pull pubsub.PullRequest
	pull.Subscription = *subname
	pull.MaxMessages = 5
	ctx := context.Background()

	for {
		resp, err := s.sub.Pull(ctx, &pull)
		if err != nil {
			fmt.Println(err)
		} else {
			s.messagesreceive(ctx, resp, pull)
		}
	}
}

func (s server) messagesreceive(ctx context.Context, resp *pubsub.PullResponse, pull pubsub.PullRequest) {
	var ackMess pubsub.AcknowledgeRequest
	ackMess.Subscription = pull.Subscription
	for _, messRec := range resp.ReceivedMessages {
		ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
		s.msgreceive(messRec.GetMessage())
	}
	s.sub.Acknowledge(ctx, &ackMess)
}

func (s server) msgreceive(msg *pubsub.PubsubMessage) {
	if starttime, err := strconv.ParseInt(msg.Attributes["time"], 10, 64); err != nil {
		fmt.Println(err)
	} else {
		var elapsedTime float64
		elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano() - starttime)
		s.timeProm.WithLabelValues("time").Observe(elapsedTime)

		var tweet twitter.Tweet
		err := json.Unmarshal(msg.Data, &tweet)
		if err != nil {
			fmt.Println(err)
		} else {
			s.tweetStream <- tweet
		}
	}
}

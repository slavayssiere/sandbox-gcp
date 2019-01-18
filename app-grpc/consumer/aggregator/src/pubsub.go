package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func (s server) connexionSubcriber(address string, filename string, scope ...string) pubsub.SubscriberClient {
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

func (s server) consumeMessageSocial() {

	var pull pubsub.PullRequest
	pull.Subscription = *subname
	pull.MaxMessages = 5
	ctx := context.Background()

	if ctx == nil {
		log.Println("Context is nil")
	}
	if s.sub == nil {
		log.Println("s.sub is nil")
	}

	for {
		resp, err := s.sub.Pull(ctx, &pull)
		if err != nil {
			fmt.Println(err)
		} else {
			s.messageSocialreceive(ctx, resp, pull)
		}
	}
}

func (s server) messageSocialreceive(ctx context.Context, resp *pubsub.PullResponse, pull pubsub.PullRequest) {
	var ackMess pubsub.AcknowledgeRequest
	ackMess.Subscription = pull.Subscription
	for _, messRec := range resp.ReceivedMessages {
		ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
		s.mSreceive(messRec.GetMessage())
	}
	s.sub.Acknowledge(ctx, &ackMess)
}

func (s server) mSreceive(msg *pubsub.PubsubMessage) {
	normtime, errn := strconv.ParseInt(msg.Attributes["normalizer_time"], 10, 64)
	if errn == nil {
		var elapsedTime float64
		elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano() - normtime)
		s.timeProm.WithLabelValues(*subname).Observe(elapsedTime)
	}

	injectime, erri := strconv.ParseInt(msg.Attributes["normalizer_time"], 10, 64)
	if erri == nil {
		var elapsedTime float64
		elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano() - injectime)
		s.timeProm.WithLabelValues(*subname).Observe(elapsedTime)
	}
	var ms libmetier.MessageSocial
	err := json.Unmarshal(msg.Data, &ms)
	if err != nil {
		log.Println(err)
	} else {
		s.messages <- (func() (libmetier.MessageSocial, int64, int64) { return ms, normtime, injectime })
	}
}

// MessageAggregator a message send by cloud sheduler
type MessageAggregator struct {
	Script string `json:"run"`
}

func (s server) consumeAggregatorMsg() {

	var pull pubsub.PullRequest
	pull.Subscription = *aggregasub
	pull.MaxMessages = 5
	ctx := context.Background()

	if ctx == nil {
		log.Println("Context is nil")
	}
	if s.sub == nil {
		log.Println("s.sub is nil")
	}

	for {
		resp, err := s.sub.Pull(ctx, &pull)
		if err != nil {
			fmt.Println(err)
		} else {
			s.messageAggregatorreceive(ctx, resp, pull)
		}
	}
}

func (s server) messageAggregatorreceive(ctx context.Context, resp *pubsub.PullResponse, pull pubsub.PullRequest) {
	var ackMess pubsub.AcknowledgeRequest
	ackMess.Subscription = pull.Subscription
	for _, messRec := range resp.ReceivedMessages {
		ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
		s.mAreceive(messRec.GetMessage())
	}
	s.sub.Acknowledge(ctx, &ackMess)
}

func (s server) mAreceive(msg *pubsub.PubsubMessage) {
	var ma MessageAggregator
	err := json.Unmarshal(msg.Data, &ma)
	if err != nil {
		log.Println(err)
	} else {
		if ma.Script == "dataset" {
			log.Println("dataset generator")
		} else {
			log.Println("aggrega generator")
			agg := s.computeAggregas()
			s.writeAggrega("aggregas", agg)
		}
	}
}

package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"log"

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

func (s server) consumemessage() {

	var pull pubsub.PullRequest
	pull.Subscription = *subname
	pull.MaxMessages = 5
	ctx := context.Background()

	for {
		if ctx == nil {
			log.Println("Context is nil")
		}
		if s.sub == nil {
			log.Println("s.sub is nil")
		}
		resp, err := s.sub.Pull(ctx, &pull)
		if err != nil {
			log.Println(err)
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
	var ms libmetier.MessageSocial
	err := json.Unmarshal(msg.Data, &ms)
	if err != nil {
		log.Println(err)
	} else {
		s.messages <- ms
	}
}

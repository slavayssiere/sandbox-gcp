package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func connexionPublisher(address string, filename string, scope ...string) pubsub.PublisherClient {
	var err error

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
	}

	creds := credentials.NewClientTLSFromCert(pool, "")

	perRPC, err := oauth.NewServiceAccountFromFile(filename, scope...)
	if err != nil {
		log.Println(err)
	}

	conn, err := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
	if err != nil {
		log.Println(err)
	}

	return pubsub.NewPublisherClient(conn)
}

func connexionSubcriber(address string, filename string, scope ...string) pubsub.SubscriberClient {
	var err error

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
	}

	creds := credentials.NewClientTLSFromCert(pool, "")

	perRPC, err := oauth.NewServiceAccountFromFile(filename, scope...)
	if err != nil {
		log.Println(err)
	}

	conn, err := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
	if err != nil {
		log.Println(err)
	}

	return pubsub.NewSubscriberClient(conn)
}

func (s server) publishmessage(mess *pubsub.PubsubMessage) {

	var request pubsub.PublishRequest
	request.Topic = *topic
	request.Messages = append(request.Messages, mess)

	fmt.Println(request)

	ctx := context.Background()

	if _, err := s.pub.Publish(ctx, &request); err != nil {
		fmt.Println(err)
	}
}

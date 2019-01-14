package main

import (
	"context"
	"io/ioutil"
	"log"
	"cloud.google.com/go/datastore"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func datastoreClient(ctx context.Context) *datastore.Client {
	jsonKey, err := ioutil.ReadFile(*secretpath)
	config, err := google.JWTConfigFromJSON(jsonKey, bigtable.Scope) // or bigtable.AdminScope, etc.
	client, err := datastore.NewClient(ctx, *projectid, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the kind for the new entity.
	kind := "test"
	// Sets the name/ID for the new entity.
	name := "ms"
	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, name, nil)

	return client
}

func (s server) writeMessage(ctx context.Context, mess libmetier.MessageSocial) {
	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &mess); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}
}

func (s server) readMessage(ctx context.Context) libmetier.MessageSocial {

}

func (s server) writeMessages(ctx context.Context) {

	for {
		mess := <-s.messages
		s.writeMessage(ctx, mess)
	}
}

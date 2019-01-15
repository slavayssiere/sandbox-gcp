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
	config, err := google.JWTConfigFromJSON(jsonKey, datastore.ScopeDatastore) // or bigtable.AdminScope, etc.
	client, err := datastore.NewClient(ctx, *projectid, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func (s server) writeMessage(ctx context.Context, mess libmetier.MessageSocial) {
	// Saves the new entity.
	k := datastore.IncompleteKey("ms", nil)
	if _, err := s.ds.Put(ctx, k, &mess); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}
}

func (s server) readMessage(ctx context.Context) *libmetier.MessageSocial {
	// mss := make([]*libmetier.MessageSocial, 0)
	// q := datastore.NewQuery("ms")

	//keys, err := s.ds.GetAll(ctx, q, &mss)

	// if err != nil {
	// 	log.Printf("datastoredb: could not list books: %v", err)
	// 	return nil
	// }

	return nil
}

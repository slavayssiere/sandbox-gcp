package main

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type twitterClient struct {
	Client *twitter.Client
	Filter []string
}

func newTwitter(consumerKey, consumerSecret, accessToken, accessSecret *string) twitterClient {
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	return twitterClient{Client: twitter.NewClient(httpClient)}
}

func (twitterClient twitterClient) filterTwitter(hashtag string) *twitter.Stream {
	log.Println("Starting Stream...")

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{hashtag},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := twitterClient.Client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	return stream
}

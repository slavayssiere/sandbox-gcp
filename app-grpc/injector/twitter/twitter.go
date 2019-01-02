package main

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterClient struct {
	Client *twitter.Client
	Filter []string
}

func NewTwitter(consumerKey, consumerSecret, accessToken, accessSecret *string) TwitterClient {
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	return TwitterClient{Client: twitter.NewClient(httpClient)}
}

func (twitterClient TwitterClient) FilterTwitter() *twitter.Stream {
	log.Println("Starting Stream...")

	filterParams := &twitter.StreamFilterParams{
		Track:         twitterClient.Filter,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := twitterClient.Client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	return stream
}

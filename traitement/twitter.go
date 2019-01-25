package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"cloud.google.com/go/bigquery"
	// libmetier "github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	"google.golang.org/api/iterator"
)

func newTwitter(consumerKey string, consumerSecret string, accessToken string, accessSecret string) *twitter.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	return twitter.NewClient(httpClient)
}

func getResults(start int64, limit int64, test *bigquery.RowIterator, wg *sync.WaitGroup, nb int64) {
	var tot Result
	test.StartIndex = uint64(start)
	var i int64
	i = start
	for {
		var row []bigquery.Value
		err := test.Next(&row)
		if err == iterator.Done {
			log.Println("finish")
			break
		}
		if err != nil {
			log.Println(err)
		}

		if row[2] == "twitter" {
			if row[0] != nil {
				id, err := strconv.ParseInt(row[0].(string), 10, 64)
				if err != nil {
					log.Println(err)
				}

				tweet, _, err := tc.Statuses.Show(id, nil)
				if err != nil {
					log.Println(err)
				} else {
					tot.FavoriteCount += tweet.FavoriteCount
					tot.RetweetCount += tweet.RetweetCount
				}
			}
		} else {
			log.Println("mastodon value")
		}
		i++
		if i >= limit {
			break
		}
	}

	log.Println(nb, ".RetweetCount:", tot.RetweetCount)
	log.Println(nb, ".FavoriteCount:", tot.FavoriteCount)

	log.Println("send result !")
	resultRetweet <- tot.RetweetCount
	resultLike <- tot.FavoriteCount

	(*wg).Done()
}

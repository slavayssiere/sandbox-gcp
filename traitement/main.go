package main

import (
	"context"
	"log"
	"os"
	"sync"

	//"strconv"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/bigtable"
	"github.com/dghubble/go-twitter/twitter"
	// libmetier "github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	//"google.golang.org/api/iterator"
)

var (
	tableName  = "ms"
	familyName = "test-table"

	// Client is initialized by main.
	client        *bigtable.Client
	tc            *twitter.Client
	resultRetweet chan int
	resultLike    chan int
)

// Result test
type Result struct {
	FavoriteCount int   `json:"nb_like"`
	RetweetCount  int   `json:"nb_retweet"`
	ID            int64 `json:"id"`
}

func main() {

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, "slavayssiere-sandbox")
	if err != nil {
		log.Println(err)
	}

	resultRetweet = make(chan int)
	resultLike = make(chan int)

	tc = newTwitter(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))

	query := client.Query("SELECT ms.iD.cell.value as id, ms.tag.cell.value as tag, ms.source.cell.value as src FROM `slavayssiere-sandbox.test_bq.ms`") //ORDER BY ms.id.cell.timestamp DESC
	iter, err := query.Read(ctx)
	if err != nil {
		log.Println(err)
	}

	var nb int64
	var block int64
	block = 1000
	nb = int64(iter.TotalRows) / int64(block)
	log.Println("parralellism:", nb, " & total: ", iter.TotalRows)

	var wg sync.WaitGroup
	wg.Add(int(nb + 1))
	var i int64
	for i = 0; i < (nb + 1); i++ {
		log.Println("launch :", i, " | from :", (i * block), " & to: ", (i+1)*block)
		go getResults(i*block, (i+1)*block, iter, &wg, i)
	}

	var totRetweetCount int
	var totFavoriteCount int
	go func() {
		for {
			ret := <-resultRetweet
			totRetweetCount += ret
			log.Println("tot.RetweetCount", totRetweetCount)
		}
	}()

	go func() {
		for {
			ret := <-resultLike
			totFavoriteCount += ret
			log.Println("tot.FavoriteCount", totFavoriteCount)
		}
	}()

	wg.Wait()
	log.Println("totRetweetCount:", totRetweetCount, " & totFavoriteCount: ", totFavoriteCount)
}

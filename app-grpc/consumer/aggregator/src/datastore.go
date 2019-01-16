package main

import (
	"context"
	"io/ioutil"
	"log"
	"sort"

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

func (s server) writeMessage(ad libmetier.AggregatedData) {
	// Saves the new entity.
	k := datastore.IncompleteKey("userstats", nil)
	if _, err := s.ds.Put(s.ctx, k, &ad); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}
}


func (s server) writeBulkMessage(ads []libmetier.AggregatedData) {
	// Saves the new entity.
	kl := make([]*datastore.Key, len(ads))
	k := datastore.IncompleteKey("userstats", nil)
	for i:=0; i!=len(ads); i++ {
		kl[i]=k
	}
	if len(ads) > 500 {
		var min int
		var max int
		min=0
		max=400
		for {
			log.Printf("min: %d, &max: %d, &len(ads): %d", min, max, len(ads))
			adstmp := ads[min:max]
			kltmp := kl[min:max]
			if _, err := s.ds.PutMulti(s.ctx, kltmp, adstmp); err != nil {
				log.Fatalf("Failed to save task: %v", err)
			}
			min = max+1
			if max > (len(ads) - 400) {
				max = len(ads)
			} else {
				max = max + 400
			}
			if min > len(ads) {
				break
			}
		}
	} else {
		if _, err := s.ds.PutMulti(s.ctx, kl, ads); err != nil {
			log.Fatalf("Failed to save task: %v", err)
		}
	}
}

func (s server) readMessage() []libmetier.AggregatedData {
	var uss []libmetier.AggregatedData
	
	q := datastore.NewQuery("userstats").Order("count").Limit(10)
	_, err := s.ds.GetAll(s.ctx, q, &uss)
	if err != nil {
		log.Printf("datastoredb: could not list books: %v", err)
		return nil
	}

	return uss
}

func (s server) createStats() []libmetier.AggregatedData {
	ads := s.getUsersCounter(-1)
	log.Println("get users, done !")
	s.writeBulkMessage(ads)
	log.Println("write users, done !")
	return ads
}

func top10(ads []libmetier.AggregatedData) []libmetier.AggregatedData{
	count := func(p1, p2 *libmetier.AggregatedData) bool {
		return p1.Count > p2.Count
	}

	By(count).Sort(ads)

	if len(ads) > 10 {
		ads=ads[:10]
	}
	return ads
}

// By is the type of a "less" function that defines the ordering of users
type By func(ad1, ad2 *libmetier.AggregatedData) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(ads []libmetier.AggregatedData) {
	as := &adSorter{
		ads: ads,
		by:  by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(as)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type adSorter struct {
	ads []libmetier.AggregatedData
	by  func(p1, p2 *libmetier.AggregatedData) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *adSorter) Len() int {
	return len(s.ads)
}

// Swap is part of sort.Interface.
func (s *adSorter) Swap(i, j int) {
	s.ads[i], s.ads[j] = s.ads[j], s.ads[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *adSorter) Less(i, j int) bool {
	return s.by(&s.ads[i], &s.ads[j])
}

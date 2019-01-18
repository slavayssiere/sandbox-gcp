package main

import (
	"sort"
	"time"

	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
)

// Aggrega test
type Aggrega struct {
	InjectorMean   float64   `json:"mean_time_injector" datastore:"mt_inj"`
	InjectorNb     int64     `json:"count_injector" datastore:"nb_inj"`
	NormalizerMean float64   `json:"mean_time_normalizer" datastore:"mt_nor"`
	NormalizerNb   int64     `json:"count_normalizer" datastore:"nb_nor"`
	Num            int64     `json:"id" datastore:"num"`
	CreateTime     time.Time `json:"create" datastore:"create_timestamp"`
}

func (s server) top10() []libmetier.AggregatedData {
	count := func(p1, p2 *libmetier.AggregatedData) bool {
		return p1.Count > p2.Count
	}

	ads := s.getUsersCounter(-1)

	By(count).Sort(ads)

	if len(ads) > 10 {
		ads = ads[:10]
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

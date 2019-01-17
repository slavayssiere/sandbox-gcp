package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s server) handlerUsersFunc(w http.ResponseWriter, r *http.Request) {
	us := s.getUsersCounter(100)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

func (s server) handlerTopTenFunc(w http.ResponseWriter, r *http.Request) {
	us := s.top10()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

type statStatus struct {
	Status  string  `json:"status" default:"done"`
	Elapsed int64   `json:"time"`
	Agg     Aggrega `json:"result"`
	Nb      int64   `json:"num_aggrega"`
}

func (s server) handlerStatsFunc(w http.ResponseWriter, r *http.Request) {
	var ret statStatus

	start := time.Now()
	ret.Agg, ret.Nb = s.computeAggregas()
	s.writeAggrega("aggregas", ret.Agg)
	t := time.Now()
	ret.Elapsed = int64(t.Sub(start))
	ret.Status = "done"

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

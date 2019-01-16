package main

import (
	"encoding/json"
	"net/http"
)

func (s server) handlerUsersFunc(w http.ResponseWriter, r *http.Request) {
	us := s.getUsersCounter(100)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

func (s server) handlerPostUsersFunc(w http.ResponseWriter, r *http.Request) {
	us := s.createStats()
	us = top10(us)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

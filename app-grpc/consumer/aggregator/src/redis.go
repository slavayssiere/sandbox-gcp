package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
)

func redisNew() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     *redisaddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

func (s server) countUser(user string) {
	s.redis.Incr(user)
	s.redis.SAdd("list_users", user)
}

func (s server) getUsersCounter() []libmetier.AggregatedData {
	var ret []libmetier.AggregatedData
	users, err := s.redis.SMembers("list_users").Result()
	if err != nil {
		log.Println(err)
	}
	for id := range users {
		var user libmetier.AggregatedData
		user.User = users[id]
		user.Count, err = s.redis.Get(users[id]).Int64()
		if err != nil {
			log.Println(err)
		}
		user.Date = time.Now()
	}
	return ret
}

func (s server) writeMessages(ctx context.Context) {

	for {
		mess := <-s.messages
		if len(mess.User) > 0 {
			s.countUser(mess.User)
		}
	}
}

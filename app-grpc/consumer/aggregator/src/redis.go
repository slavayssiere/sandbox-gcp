package main

import (
	"log"
	"strconv"
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

func (s server) getUsersCounter(limit int) []libmetier.AggregatedData {
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
		ret = append(ret, user)
		if limit > 0 {
			if id > limit {
				break
			}
		}
	}
	return ret
}

func (s server) addNormTime(normtime int64) {
	s.redis.LPush("normTimes_"+string(s.getNbAggregation()), normtime)
}

func (s server) addInjectTime(injectime int64) {
	s.redis.LPush("injectTimes_"+string(s.getNbAggregation()), injectime)
}

func (s server) addAggTime(aggtime int64) {
	s.redis.LPush("aggTimes_"+string(s.getNbAggregation()), aggtime)
}

func (s server) getMeanTimes(key string, aggrega int64) (float64, int64) {
	nb, erra := s.redis.LLen(key + string(aggrega)).Result()
	if erra != nil {
		log.Println(erra)
	}
	val, errb := s.redis.LRange(key+string(aggrega), 0, nb).Result()
	if errb != nil {
		log.Println(errb)
	}
	s.redis.Del(key + string(aggrega))
	var sum int64
	var i int64
	sum = 0
	for i = 0; i != nb; i++ {
		tmp, _ := strconv.ParseInt(val[i], 10, 64)
		sum = sum + tmp
	}
	var ret float64
	ret = (float64(sum) / float64(nb))
	return ret, nb
}

func (s server) computeAggregas() Aggrega {
	var agg Aggrega

	agg.Num = s.getNbAggregation()
	s.addAggregation()

	agg.InjectorMean, agg.InjectorNb = s.getMeanTimes("injectTimes_", agg.Num)
	agg.NormalizerMean, agg.NormalizerNb = s.getMeanTimes("normTimes_", agg.Num)

	return agg
}

func (s server) addAggregation() {
	s.redis.Incr("aggregas")
}

func (s server) getNbAggregation() int64 {
	nb, err := s.redis.Get("aggregas").Int64()
	if err != nil {
		log.Println(err)
	}
	return nb
}

func (s server) writeMessagesToRedis() {

	for {
		mess, normtime, injectime := (<-s.messages)()
		if len(mess.User) > 0 {
			s.countUser(mess.User)
		}
		s.addNormTime(time.Now().UnixNano() - normtime)
		s.addInjectTime(time.Now().UnixNano() - injectime)
	}
}

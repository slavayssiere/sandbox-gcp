package main

import "github.com/prometheus/client_golang/prometheus"

func PromHistogramVec() *prometheus.HistogramVec {
	histogramMean := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "mean_in_injector",
			Help: "Time for pubish to pubsub in nanosecond",
		},
		[]string{"size", "trade"},
	)

	prometheus.Register(histogramMean)

	return histogramMean
}

func PromCounterVec() *prometheus.CounterVec {
	messagesCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_injected",
			Help: "How many messages injected, partitioned by size and trade",
		},
		[]string{"size", "trade"},
	)

	prometheus.Register(messagesCounter)

	return messagesCounter
}

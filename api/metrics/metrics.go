package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "lambda_api_requests_total",
		Help: "Total number of requests",
	}, []string{"path", "method", "status"})

	RequestsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "lambda_api_requests_duration",
		Help:    "Duration of requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"path", "method", "status"})

	PartyCountTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "lambda_api_party_count_total",
		Help: "Total number of parties",
	})

	SuccessfulLogins = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lambda_api_successful_logins",
		Help: "Total number of successful logins",
	})

	FailedLogins = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lambda_api_failed_logins",
		Help: "Total number of failed logins",
	})
)

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsCounter       prometheus.Counter
	ResponseCounter       *prometheus.CounterVec
	HistogramResponseTime *prometheus.HistogramVec
)

func Init(namespace string) {
	RequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "grpc",
		Name:      "requests_total",
	})
	ResponseCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "grpc",
		Name:      "responses_total",
	},
		[]string{"status"},
	)
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status"},
	)
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CacheHitCount      prometheus.Counter
	CacheRequestsTotal prometheus.Counter
	CacheMissCount     prometheus.Counter
)

func Init(namespace string) {
	CacheHitCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "cache_hits_total",
	})

	CacheRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "cache_requests_total",
	})

	CacheMissCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "cache_miss_total",
	})
}

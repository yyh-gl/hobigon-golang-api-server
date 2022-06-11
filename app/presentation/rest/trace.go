package rest

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	httpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "hobigon",
			Name:      "http_requests_total",
			Help:      "The number of requests",
		},
		[]string{"method", "path", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "hobigon",
			Name:      "http_response_time_seconds",
			Help:      "Duration of HTTP requests.",
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestTotal, httpDuration, collectors.NewBuildInfoCollector())
}

// IncrementRequestCount : httpRequestTotalカウンターをインクリメント
func IncrementRequestCount(method, path string, statusCode int) {
	counter := httpRequestTotal.WithLabelValues(method, path, strconv.Itoa(statusCode))
	counter.Inc()
}

// ObserveLatency : レイテンシを測定
func ObserveLatency(method, path string) func() time.Duration {
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(method, path))
	return timer.ObserveDuration
}

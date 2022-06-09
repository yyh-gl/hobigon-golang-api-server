package rest

import "github.com/prometheus/client_golang/prometheus"

var httpReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "hobigon",
		Name:      "request_count",
		Help:      "The number of requests",
	},
	[]string{"method", "path"},
	// TODO: ステータスコードを記録
	//[]string{"method", "path", "status_code"},
)

func init() {
	prometheus.MustRegister(httpReqs)
}

// CountRequest : リクエストをカウントアップ
func CountRequest(method, path string) {
	counter := httpReqs.WithLabelValues(method, path)
	counter.Inc()
}

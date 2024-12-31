package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	inFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"handler", "code", "method", "path"},
	)

	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)

	runningVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "running_version",
		Help: "A gauge of running version.",
	}, []string{"version"})

	requestCountPerPath = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "request_count_per_path",
		Help: "A gauge of request count per path.",
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(inFlight, counter, duration, responseSize, runningVersion, requestCountPerPath)
}

func InstrumentPrometheus(h http.HandlerFunc, path Path, handlerName string) http.HandlerFunc {
	return promhttp.InstrumentHandlerInFlight(inFlight,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": handlerName}),
			promhttp.InstrumentHandlerCounter(counter.MustCurryWith(prometheus.Labels{"handler": handlerName, "path": path}),
				promhttp.InstrumentHandlerResponseSize(responseSize, h),
			),
		),
	).(http.HandlerFunc)
}

func CountUpRunningVersion(version string) {
	runningVersion.With(prometheus.Labels{"version": version}).Inc()
}

func CountDownRunningVersion(version string) {
	runningVersion.With(prometheus.Labels{"version": version}).Dec()
}

func CountUpRequestCountPerPath(path string) {
	requestCountPerPath.With(prometheus.Labels{"path": path}).Inc()
}

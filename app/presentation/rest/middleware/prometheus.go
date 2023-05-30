package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	InFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	Counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"handler", "code", "method"},
	)

	Duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method"},
	)

	ResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)

	RunningVersion = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "running_version",
		Help: "A gauge of running version.",
	}, []string{"version"})
)

//func init() {
//	prometheus.MustRegister(inFlight, counter, duration, responseSize, runningVersion)
//}

func prometheusInstrument(h http.HandlerFunc, name string) http.HandlerFunc {

	// TODO: https://christina04.hatenablog.com/entry/prometheus-application-over-http 全部やりきる

	return promhttp.InstrumentHandlerInFlight(InFlight,
		promhttp.InstrumentHandlerDuration(Duration.MustCurryWith(prometheus.Labels{"handler": name}),
			promhttp.InstrumentHandlerCounter(Counter.MustCurryWith(prometheus.Labels{"handler": name}),
				promhttp.InstrumentHandlerResponseSize(ResponseSize, h),
			),
		),
	).(http.HandlerFunc)
}

// FIXME: ミドルウェア的な機能ではないので別の場所に移動させたい
func CountUpRunningVersion(version string) {
	RunningVersion.With(prometheus.Labels{"version": version}).Inc()
}

func CountDownRunningVersion(version string) {
	RunningVersion.With(prometheus.Labels{"version": version}).Dec()
}

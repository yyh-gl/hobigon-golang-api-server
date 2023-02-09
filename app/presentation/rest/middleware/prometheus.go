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
		[]string{"handler", "code", "method"},
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
)

func init() {
	prometheus.MustRegister(inFlight, counter, duration, responseSize, runningVersion)
}

func prometheusInstrument(h http.HandlerFunc, name string) http.HandlerFunc {

	// TODO https://christina04.hatenablog.com/entry/prometheus-application-over-http 全部やりきる

	return promhttp.InstrumentHandlerInFlight(inFlight,
		promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": name}),
			promhttp.InstrumentHandlerCounter(counter.MustCurryWith(prometheus.Labels{"handler": name}),
				promhttp.InstrumentHandlerResponseSize(responseSize, h),
			),
		),
	).(http.HandlerFunc)
}

// FIXME: ミドルウェア的な機能ではないので別の場所に移動させたい
// FIXME: 正しくカウントアップできていないので修正する
// FIXME: バージョン情報にコンテナイメージの名称が含まれているので削除
func CountUpRunningVersion(version string) {
	runningVersion.With(prometheus.Labels{"version": version}).Inc()
}

func CountDownRunningVersion(version string) {
	runningVersion.With(prometheus.Labels{"version": version}).Dec()
}

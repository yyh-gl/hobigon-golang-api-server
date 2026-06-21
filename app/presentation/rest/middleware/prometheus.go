package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	inFlight             metric.Int64UpDownCounter
	requestCounter       metric.Int64Counter
	requestDuration      metric.Float64Histogram
	responseSizeHist     metric.Int64Histogram
	runningVersionGauge  metric.Int64UpDownCounter
	requestCountPerPath  metric.Int64Counter
)

func init() {
	m := otel.GetMeterProvider().Meter("hobigon-rest")

	var err error
	inFlight, err = m.Int64UpDownCounter("http.server.active_requests",
		metric.WithDescription("Number of requests currently being served."))
	if err != nil {
		panic(err)
	}

	requestCounter, err = m.Int64Counter("http.server.request.total",
		metric.WithDescription("Total number of HTTP requests."))
	if err != nil {
		panic(err)
	}

	requestDuration, err = m.Float64Histogram("http.server.request.duration",
		metric.WithDescription("Duration of HTTP requests in seconds."),
		metric.WithUnit("s"))
	if err != nil {
		panic(err)
	}

	responseSizeHist, err = m.Int64Histogram("http.server.response.size",
		metric.WithDescription("Size of HTTP responses in bytes."),
		metric.WithUnit("By"))
	if err != nil {
		panic(err)
	}

	runningVersionGauge, err = m.Int64UpDownCounter("app.running_version",
		metric.WithDescription("Currently running application version."))
	if err != nil {
		panic(err)
	}

	requestCountPerPath, err = m.Int64Counter("http.server.request.count_per_path",
		metric.WithDescription("Request count per URL path."))
	if err != nil {
		panic(err)
	}
}

func InstrumentOTel(h http.HandlerFunc, path Path, handlerName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		attrs := []attribute.KeyValue{
			attribute.String("handler", handlerName),
			attribute.String("method", r.Method),
			attribute.String("path", path),
		}

		inFlight.Add(ctx, 1, metric.WithAttributes(attrs...))
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w}
		h.ServeHTTP(rw, r)

		elapsed := time.Since(start).Seconds()
		statusAttrs := append(attrs, attribute.String("code", strconv.Itoa(rw.status)))

		inFlight.Add(ctx, -1, metric.WithAttributes(attrs...))
		requestCounter.Add(ctx, 1, metric.WithAttributes(statusAttrs...))
		requestDuration.Record(ctx, elapsed, metric.WithAttributes(attrs...))
		responseSizeHist.Record(ctx, int64(rw.size))
	}
}

func CountUpRunningVersion(version string) {
	runningVersionGauge.Add(context.Background(), 1,
		metric.WithAttributes(attribute.String("version", version)))
}

func CountDownRunningVersion(version string) {
	runningVersionGauge.Add(context.Background(), -1,
		metric.WithAttributes(attribute.String("version", version)))
}

func CountUpRequestCountPerPath(path string) {
	requestCountPerPath.Add(context.Background(), 1,
		metric.WithAttributes(attribute.String("path", path)))
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

package telemetry

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func SetupOTel(ctx context.Context) (shutdown func(context.Context) error, err error) {
	if os.Getenv("APP_ENV") == "test" {
		return func(context.Context) error { return nil }, nil
	}

	var shutdowns []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var errs []error
		for i := len(shutdowns) - 1; i >= 0; i-- {
			errs = append(errs, shutdowns[i](ctx))
		}
		return errors.Join(errs...)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("hobigon-rest"),
			semconv.ServiceVersion(app.GetVersion()),
		),
	)
	if err != nil {
		return shutdown, err
	}

	traceExporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return shutdown, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	shutdowns = append(shutdowns, tp.Shutdown)

	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return shutdown, err
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
	shutdowns = append(shutdowns, mp.Shutdown)

	logExporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return shutdown, err
	}
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
	global.SetLoggerProvider(lp)
	shutdowns = append(shutdowns, lp.Shutdown)

	return shutdown, nil
}

package log

import (
	"context"

	"go.opentelemetry.io/otel"
)

func SetTraceIDToContext(ctx context.Context) context.Context {
	ctx, span := otel.Tracer("hobigon").Start(ctx, "request")
	_ = span
	return ctx
}

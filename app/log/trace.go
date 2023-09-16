package log

import (
	"github.com/rs/xid"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"golang.org/x/net/context"
)

func NewTraceID() string {
	return xid.New().String()
}

func SetTraceIDToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, app.TraceIdContextKey, NewTraceID())
}

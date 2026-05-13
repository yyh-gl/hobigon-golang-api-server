package log

import (
	"context"

	"github.com/rs/xid"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func NewTraceID() string {
	return xid.New().String()
}

func SetTraceIDToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, app.ContextKeyTraceId, NewTraceID())
}

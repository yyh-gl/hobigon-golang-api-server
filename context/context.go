package context

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

type contextKey int

const (
	ParamsKey contextKey = iota
)

func InjectRequestParams(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, ParamsKey, params)
}

func FetchRequestParams(ctx context.Context) httprouter.Params {
	return ctx.Value(ParamsKey).(httprouter.Params)
}

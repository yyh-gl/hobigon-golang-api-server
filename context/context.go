package context

import (
	"context"

	"github.com/julienschmidt/httprouter"
)

type contextKey int

// ParamsKey : リクエストパラメータを取得するためのキー名
const ParamsKey contextKey = iota

// InjectRequestParams : リクエストパラメータを格納
func InjectRequestParams(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, ParamsKey, params)
}

// FetchRequestParams : リクエストパラメータを取得
func FetchRequestParams(ctx context.Context) httprouter.Params {
	return ctx.Value(ParamsKey).(httprouter.Params)
}

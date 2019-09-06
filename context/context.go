package context

import (
	"context"
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

type contextKey int

const (
	ParamsKey contextKey = iota
	LoggerKey
	DBKey
)

func InjectRequestParams(ctx context.Context, params httprouter.Params) context.Context {
	return context.WithValue(ctx, ParamsKey, params)
}

func FetchRequestParams(ctx context.Context) httprouter.Params {
	return ctx.Value(ParamsKey).(httprouter.Params)
}

func InjectLogger(ctx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func FetchLogger(ctx context.Context) (*log.Logger, error) {
	logger, ok := ctx.Value(LoggerKey).(*log.Logger)
	if !ok {
		return nil, errors.New("Failed to fetch logger from context")
	}
	return logger, nil
}

func InjectDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, DBKey, db)
}

func FetchDB(ctx context.Context) *gorm.DB {
	return ctx.Value(DBKey).(*gorm.DB)
}

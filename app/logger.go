package app

import (
	"context"
	"log/slog"
	"os"
)

var logger *slog.Logger

func NewLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Info(ctx context.Context, msg string) {
	logger.InfoContext(
		ctx,
		msg,
		slog.String("version", GetVersion()),
	)
}

func Error(ctx context.Context, err error) {
	logger.ErrorContext(
		ctx,
		err.Error(),
		slog.String("version", GetVersion()),
	)
}

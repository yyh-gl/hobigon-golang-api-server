package log

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

var logger *slog.Logger

func NewLogger() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Info(ctx context.Context, msg string) {
	traceID := ctx.Value(app.TraceIdContextKey)
	// The only time traceID is nil is in the server start log.
	if traceID == nil {
		traceID = ""
	}

	logger.InfoContext(
		ctx,
		msg,
		slog.String("version", app.GetVersion()),
		slog.String("trace_id", traceID.(string)),
	)
}

func InfoRequestAndResponse(ctx context.Context, req http.Request, resp ResponseRecorder) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		Error(ctx, err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	logger.InfoContext(
		ctx,
		"request and response log",
		slog.String("version", app.GetVersion()),
		slog.String("trace_id", ctx.Value(app.TraceIdContextKey).(string)),
		slog.Group("request",
			slog.String("method", req.Method),
			slog.String("host", req.Host),
			slog.String("uri", req.RequestURI),
			slog.String("user_agent", req.Header.Get("user-agent")),
			slog.String("body", string(body)),
			slog.Int64("content_length", req.ContentLength),
			slog.String("remote_addr", req.RemoteAddr),
		),
		slog.Group("response",
			slog.Int("status_code", resp.StatusCode),
			slog.String("body", resp.Body.String()),
		),
	)
}

func Error(ctx context.Context, err error) {
	logger.ErrorContext(
		ctx,
		err.Error(),
		slog.String("version", app.GetVersion()),
		slog.String("trace_id", ctx.Value(app.TraceIdContextKey).(string)))
}

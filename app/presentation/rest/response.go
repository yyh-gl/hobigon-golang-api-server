package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// DoResponse : JSON形式でレスポンスを返す
func DoResponse(ctx context.Context, w http.ResponseWriter, resp any, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		app.Error(ctx, fmt.Errorf("failed to json.NewEncoder().Encode(): %w", err))
		http.Error(w, "failed to create response", http.StatusInternalServerError)
		return
	}
}

// DoImageResponse : 画像ファイルをレスポンスとして返す
func DoImageResponse(ctx context.Context, w http.ResponseWriter, img []byte, contentType string, status int) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))
	w.WriteHeader(status)
	if _, err := w.Write(img); err != nil {
		app.Error(ctx, fmt.Errorf("failed to http.ResponseWriter.Wrire(): %w", err))
		DoResponse(ctx, w, "failed to create response", http.StatusInternalServerError)
	}
}

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// DoResponse : JSON形式でレスポンスを返す
func DoResponse(w http.ResponseWriter, resp any, status int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		app.Error(fmt.Errorf("json.NewEncoder().Encode()でエラー: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DoImageResponse : 画像ファイルをレスポンスとして返す
func DoImageResponse(w http.ResponseWriter, img []byte, contentType string, status int) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(img)))
	w.WriteHeader(status)
	if _, err := w.Write(img); err != nil {
		app.Error(fmt.Errorf("http.ResponseWriter.Wrire(): %w", err))
		DoResponse(w, "creating response is failed", http.StatusInternalServerError)
	}
}

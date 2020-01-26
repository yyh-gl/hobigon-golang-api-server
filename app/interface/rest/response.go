package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// Response : REST API 用の共通エラーレスポンス
type errorResponse struct {
	Error string `json:"error"`
}

// DoResponse : JSON形式でレスポンスを返す
func DoResponse(w http.ResponseWriter, resp interface{}, status int) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		app.Logger.Println(fmt.Errorf("json.NewEncoder().Encode()でエラー: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

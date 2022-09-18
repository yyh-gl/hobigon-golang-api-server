package middleware

import (
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// CORS : CORS関連のヘッダーを付与
func CORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS用ヘッダーを付与
		switch {
		case app.IsPrd():
			w.Header().Set("Access-Control-Allow-Origin", "https://tech.yyh-gl.dev")
		case app.IsDev() || app.IsTest():
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		h.ServeHTTP(w, r)
	}
}

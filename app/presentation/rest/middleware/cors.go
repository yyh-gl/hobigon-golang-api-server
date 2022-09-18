package middleware

import (
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// cors : CORS関連のヘッダーを付与
func cors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS用ヘッダーを付与
		switch {
		case app.IsPrd():
			w.Header().Add("Access-Control-Allow-Origin", "https://tech.yyh-gl.dev")
		case app.IsDev() || app.IsTest():
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
		}
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		h.ServeHTTP(w, r)
	}
}

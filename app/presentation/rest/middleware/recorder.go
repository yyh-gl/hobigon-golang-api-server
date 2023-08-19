package middleware

import (
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

// Recorder : リクエストおよびレスポンスを記録
func Recorder(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := log.NewResponseRecorder(w)

		h.ServeHTTP(&rec, r)

		log.InfoRequestAndResponse(r.Context(), *r, rec)
	}
}

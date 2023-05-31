package middleware

import (
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func Attach(h http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h = CORS(h)

		if app.IsPrd() {
			h = prometheusInstrument(h, name)
		}

		h.ServeHTTP(w, r)
	}
}

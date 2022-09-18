package middleware

import (
	"net/http"
)

func Attach(h http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h = cors(h)
		h = prometheusInstrument(h, name)

		h.ServeHTTP(w, r)
	}
}

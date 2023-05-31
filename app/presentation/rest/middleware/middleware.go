package middleware

import "net/http"

func Attach(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h = CORS(h)

		h.ServeHTTP(w, r)
	}
}

package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func Attach(h http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h = CORS(h)

		Counter.With(prometheus.Labels{"handler": "test", "code": "200", "method": "GET"}).Inc()

		if app.IsPrd() {
			h = prometheusInstrument(h, name)
		}

		h.ServeHTTP(w, r)
	}
}

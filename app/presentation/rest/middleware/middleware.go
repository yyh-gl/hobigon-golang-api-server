package middleware

import (
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

type Path = string

func CreateHandlerFuncWithMiddleware(h http.HandlerFunc, path Path, handlerName string) (Path, http.HandlerFunc) {
	h = CORS(h)
	h = Recorder(h)
	h = InstrumentPrometheus(h, path, handlerName)

	return path, func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(log.SetTraceIDToContext(r.Context()))
		h.ServeHTTP(w, r)
	}
}

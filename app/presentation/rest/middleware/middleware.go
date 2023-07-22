package middleware

import "net/http"

type Path = string

func CreateHandlerFuncWithMiddleware(h http.HandlerFunc, path Path, handlerName string) (Path, http.HandlerFunc) {
	h = CORS(h)
	h = InstrumentPrometheus(h, path, handlerName)

	return path, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

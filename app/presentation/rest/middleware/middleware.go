package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Path = string

func CreateHandlerFuncWithMiddleware(h http.HandlerFunc, path Path, handlerName string) (Path, http.HandlerFunc) {
	h = CORS(h)
	h = Recorder(h)
	h = InstrumentOTel(h, path, handlerName)

	return path, otelhttp.NewHandler(h, handlerName).ServeHTTP
}

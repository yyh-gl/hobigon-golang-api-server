package log

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) ResponseRecorder {
	return ResponseRecorder{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.StatusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *ResponseRecorder) Write(data []byte) (int, error) {
	rec.Body.Write(data)
	return rec.ResponseWriter.Write(data)
}

func (rec *ResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rec.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter doesn't satisfy Hijacker interface")
	}
	return hijacker.Hijack()
}

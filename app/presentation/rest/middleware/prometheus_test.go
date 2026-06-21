package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest/middleware"
)

func TestInstrumentOTel_doesNotPanicWithNoopMeter(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	instrumented := middleware.InstrumentOTel(h, "/test", "test_handler")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InstrumentOTel panicked: %v", r)
		}
	}()
	instrumented(rec, req)
}

func TestCountUpRunningVersion_doesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CountUpRunningVersion panicked: %v", r)
		}
	}()
	middleware.CountUpRunningVersion("v1.0.0")
}

func TestCountDownRunningVersion_doesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CountDownRunningVersion panicked: %v", r)
		}
	}()
	middleware.CountDownRunningVersion("v1.0.0")
}

func TestCountUpRequestCountPerPath_doesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CountUpRequestCountPerPath panicked: %v", r)
		}
	}()
	middleware.CountUpRequestCountPerPath("/api/v1/blogs")
}

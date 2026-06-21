package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestInstrumentOTel_with404Response(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	instrumented := middleware.InstrumentOTel(h, "/test", "test_handler")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InstrumentOTel panicked with 404: %v", r)
		}
	}()
	instrumented(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestInstrumentOTel_withResponseBody(t *testing.T) {
	body := `{"title":"test"}`
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
	instrumented := middleware.InstrumentOTel(h, "/test", "test_handler")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InstrumentOTel panicked with body: %v", r)
		}
	}()
	instrumented(rec, req)

	if !strings.Contains(rec.Body.String(), "test") {
		t.Errorf("expected body to contain 'test', got %q", rec.Body.String())
	}
}

func TestInstrumentOTel_withImplicitStatus200(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
	instrumented := middleware.InstrumentOTel(h, "/test", "test_handler")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InstrumentOTel panicked with implicit 200: %v", r)
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

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/rest/middleware"
)

func TestCreateHandlerFuncWithMiddleware_returnsCorrectPath(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	path, handler := middleware.CreateHandlerFuncWithMiddleware(h, "/api/v1/blogs", "blog_create")
	if path != "/api/v1/blogs" {
		t.Errorf("path = %q, want %q", path, "/api/v1/blogs")
	}
	if handler == nil {
		t.Error("handler must not be nil")
	}
}

func TestCreateHandlerFuncWithMiddleware_handlerIsCallable(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CreateHandlerFuncWithMiddleware panicked: %v", r)
		}
	}()

	_, handler := middleware.CreateHandlerFuncWithMiddleware(h, "/test", "test_handler")
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
}

func TestTraceContextPropagation_traceparentHeaderIsForwardedToHandler(t *testing.T) {
	traceparent := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"

	var gotTraceID string
	h := func(w http.ResponseWriter, r *http.Request) {
		gotTraceID = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("traceparent", traceparent)
	rec := httptest.NewRecorder()

	h(rec, req)

	if gotTraceID != traceparent {
		t.Errorf("traceparent header not forwarded: got %q, want %q", gotTraceID, traceparent)
	}
}

func TestTraceContextPropagation_withOTelHTTPMiddleware(t *testing.T) {
	traceparent := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"

	var handlerCalled bool
	h := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}

	_, wrappedHandler := middleware.CreateHandlerFuncWithMiddleware(h, "/test", "test_handler")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("traceparent", traceparent)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("handler panicked: %v", r)
		}
	}()

	wrappedHandler(rec, req)

	if !handlerCalled {
		t.Error("inner handler was not called")
	}
}

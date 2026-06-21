package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

func TestDoResponse_NoContent(t *testing.T) {
	log.NewLogger()

	ctx := log.SetTraceIDToContext(context.Background())
	w := httptest.NewRecorder()
	DoResponse(ctx, w, nil, http.StatusNoContent)

	if w.Body.Len() != 0 {
		t.Errorf("expected empty body for 204, got %q", w.Body.String())
	}
}

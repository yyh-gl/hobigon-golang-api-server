package telemetry_test

import (
	"context"
	"os"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/cmd/rest/telemetry"
)

func TestSetupOTel_TestEnvReturnsNoopShutdown(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	shutdown, err := telemetry.SetupOTel(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shutdown == nil {
		t.Fatal("shutdown function must not be nil")
	}

	// no-opなのでCollector接続なしにshutdownできる
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown returned error: %v", err)
	}
}

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

func TestSetupOTel_TestEnvShutdownIsIdempotent(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	shutdown, err := telemetry.SetupOTel(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 複数回shutdownしてもpanicしない
	for i := 0; i < 3; i++ {
		if err := shutdown(context.Background()); err != nil {
			t.Fatalf("shutdown[%d] returned error: %v", i, err)
		}
	}
}

func TestSetupOTel_TestEnvWithCancelledContext(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	shutdown, err := telemetry.SetupOTel(ctx)
	if err != nil {
		t.Fatalf("unexpected error with cancelled context: %v", err)
	}
	if shutdown == nil {
		t.Fatal("shutdown function must not be nil")
	}

	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown returned error: %v", err)
	}
}

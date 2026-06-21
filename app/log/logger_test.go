package log_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

func TestNewLogger_doesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewLogger panicked: %v", r)
		}
	}()
	log.NewLogger()
}

func TestInfo_doesNotPanic(t *testing.T) {
	log.NewLogger()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Info panicked: %v", r)
		}
	}()
	log.Info(context.Background(), "test message")
}

func TestError_doesNotPanic(t *testing.T) {
	log.NewLogger()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Error panicked: %v", r)
		}
	}()
	log.Error(context.Background(), errors.New("test error"))
}

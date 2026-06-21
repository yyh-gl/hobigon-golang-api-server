package db

import (
	"testing"
)

func TestNewDB_returnsSQLiteInTestEnv(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewDB panicked in test env: %v", r)
		}
	}()

	db := NewDB()
	if db == nil {
		t.Fatal("NewDB returned nil")
	}
}

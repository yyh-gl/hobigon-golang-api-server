package line_test

import (
	"errors"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

func TestLINEMessages_FindMessageFor(t *testing.T) {
	messages := line.LINEMessages{"coop": "今日は生協に入金する日"}

	got, err := messages.FindMessageFor("coop")
	if err != nil {
		t.Errorf("FindMessageFor() error = %v, want nil", err)
	}
	if got != "今日は生協に入金する日" {
		t.Errorf("FindMessageFor() = %q, want %q", got, "今日は生協に入金する日")
	}
}

func TestLINEMessages_FindMessageFor_KeyNotFound(t *testing.T) {
	messages := line.LINEMessages{"coop": "今日は生協に入金する日"}

	got, err := messages.FindMessageFor("unknown")
	if !errors.Is(err, line.ErrLINEMessageKeyNotFound) {
		t.Errorf("FindMessageFor() error = %v, want %v", err, line.ErrLINEMessageKeyNotFound)
	}
	if got != "" {
		t.Errorf("FindMessageFor() = %q, want empty string", got)
	}
}

func TestLINEMessages_FindMessageFor_MultipleKeys(t *testing.T) {
	messages := line.LINEMessages{
		"coop":      "今日は生協に入金する日",
		"seisenkan": "今日は生鮮館に入金する日",
	}

	got, err := messages.FindMessageFor("seisenkan")
	if err != nil {
		t.Errorf("FindMessageFor() error = %v, want nil", err)
	}
	if got != "今日は生鮮館に入金する日" {
		t.Errorf("FindMessageFor() = %q, want %q", got, "今日は生鮮館に入金する日")
	}
}

func TestLINEMessages_FindMessageFor_EmptyStringKey(t *testing.T) {
	messages := line.LINEMessages{"coop": "今日は生協に入金する日"}

	got, err := messages.FindMessageFor("")
	if !errors.Is(err, line.ErrLINEMessageKeyNotFound) {
		t.Errorf("FindMessageFor() error = %v, want %v", err, line.ErrLINEMessageKeyNotFound)
	}
	if got != "" {
		t.Errorf("FindMessageFor() = %q, want empty string", got)
	}
}

func TestLINEMessages_FindMessageFor_EmptyMap(t *testing.T) {
	messages := line.LINEMessages{}

	got, err := messages.FindMessageFor("coop")
	if !errors.Is(err, line.ErrLINEMessageKeyNotFound) {
		t.Errorf("FindMessageFor() error = %v, want %v", err, line.ErrLINEMessageKeyNotFound)
	}
	if got != "" {
		t.Errorf("FindMessageFor() = %q, want empty string", got)
	}
}

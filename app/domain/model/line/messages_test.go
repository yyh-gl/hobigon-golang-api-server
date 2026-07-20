package line_test

import (
	"errors"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

func TestLINEMessages_MessageFor(t *testing.T) {
	messages := line.LINEMessages{"coop": "今日は生協に入金する日"}

	got, err := messages.MessageFor("coop")
	if err != nil {
		t.Errorf("MessageFor() error = %v, want nil", err)
	}
	if got != "今日は生協に入金する日" {
		t.Errorf("MessageFor() = %q, want %q", got, "今日は生協に入金する日")
	}
}

func TestLINEMessages_MessageFor_KeyNotFound(t *testing.T) {
	messages := line.LINEMessages{"coop": "今日は生協に入金する日"}

	got, err := messages.MessageFor("unknown")
	if !errors.Is(err, line.ErrLINEMessageKeyNotFound) {
		t.Errorf("MessageFor() error = %v, want %v", err, line.ErrLINEMessageKeyNotFound)
	}
	if got != "" {
		t.Errorf("MessageFor() = %q, want empty string", got)
	}
}

func TestLINEMessages_MessageFor_MultipleKeys(t *testing.T) {
	messages := line.LINEMessages{
		"coop":      "今日は生協に入金する日",
		"seisenkan": "今日は生鮮館に入金する日",
	}

	got, err := messages.MessageFor("seisenkan")
	if err != nil {
		t.Errorf("MessageFor() error = %v, want nil", err)
	}
	if got != "今日は生鮮館に入金する日" {
		t.Errorf("MessageFor() = %q, want %q", got, "今日は生鮮館に入金する日")
	}
}

func TestLINEMessages_MessageFor_EmptyStringKey(t *testing.T) {
	messages := line.LINEMessages{"coop": "今日は生協に入金する日"}

	got, err := messages.MessageFor("")
	if !errors.Is(err, line.ErrLINEMessageKeyNotFound) {
		t.Errorf("MessageFor() error = %v, want %v", err, line.ErrLINEMessageKeyNotFound)
	}
	if got != "" {
		t.Errorf("MessageFor() = %q, want empty string", got)
	}
}

func TestLINEMessages_MessageFor_EmptyMap(t *testing.T) {
	messages := line.LINEMessages{}

	got, err := messages.MessageFor("coop")
	if !errors.Is(err, line.ErrLINEMessageKeyNotFound) {
		t.Errorf("MessageFor() error = %v, want %v", err, line.ErrLINEMessageKeyNotFound)
	}
	if got != "" {
		t.Errorf("MessageFor() = %q, want empty string", got)
	}
}

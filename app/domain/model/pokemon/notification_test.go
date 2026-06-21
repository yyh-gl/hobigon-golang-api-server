package pokemon_test

import (
	"testing"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/pokemon"
)

func TestIsEventCategory(t *testing.T) {
	tests := []struct {
		category string
		want     bool
	}{
		{"イベント", true},
		{"大会", false},
		{"", false},
	}
	for _, tt := range tests {
		n := pokemon.NewNotification(tt.category, "title", "2024.1.1")
		if got := n.IsEventCategory(); got != tt.want {
			t.Errorf("IsEventCategory(%q) = %v, want %v", tt.category, got, tt.want)
		}
	}
}

func TestIsImportantEvent(t *testing.T) {
	tests := []struct {
		title string
		want  bool
	}{
		{"シティリーグ予選", true},
		{"シティリーグ", true},
		{"通常大会", false},
		{"", false},
	}
	for _, tt := range tests {
		n := pokemon.NewNotification("イベント", tt.title, "2024.1.1")
		if got := n.IsImportantEvent(); got != tt.want {
			t.Errorf("IsImportantEvent(%q) = %v, want %v", tt.title, got, tt.want)
		}
	}
}

func TestIsReceivedInToday(t *testing.T) {
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		date string
		want bool
	}{
		{"2024.6.15", true},
		{"2024.6.14", false},
		{"2024.6.16", false},
	}
	for _, tt := range tests {
		n := pokemon.NewNotification("イベント", "title", tt.date)
		if got := n.IsReceivedInToday(now); got != tt.want {
			t.Errorf("IsReceivedInToday(%q) = %v, want %v", tt.date, got, tt.want)
		}
	}
}

func TestIsReceivedInYesterday(t *testing.T) {
	now := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		date string
		want bool
	}{
		{"2024.6.14", true},
		{"2024.6.15", false},
		{"2024.6.13", false},
	}
	for _, tt := range tests {
		n := pokemon.NewNotification("イベント", "title", tt.date)
		if got := n.IsReceivedInYesterday(now); got != tt.want {
			t.Errorf("IsReceivedInYesterday(%q) = %v, want %v", tt.date, got, tt.want)
		}
	}
}

func TestTitle(t *testing.T) {
	n := pokemon.NewNotification("イベント", "シティリーグ", "2024.6.15")
	if got := n.Title(); got != "シティリーグ" {
		t.Errorf("Title() = %q, want \"シティリーグ\"", got)
	}
}

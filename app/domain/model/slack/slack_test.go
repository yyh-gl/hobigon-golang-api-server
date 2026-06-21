package slack_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/slack"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

func TestGetWebHookURL(t *testing.T) {
	tests := []struct {
		channel string
		envKey  string
		envVal  string
		want    string
	}{
		{"00_today_tasks", "WEBHOOK_URL_TO_00", "https://hooks.slack.com/00", "https://hooks.slack.com/00"},
		{"03_pokemon", "WEBHOOK_URL_TO_03", "https://hooks.slack.com/03", "https://hooks.slack.com/03"},
		{"51_tech_blog", "WEBHOOK_URL_TO_51", "https://hooks.slack.com/51", "https://hooks.slack.com/51"},
		{"unknown", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.channel, func(t *testing.T) {
			if tt.envKey != "" {
				t.Setenv(tt.envKey, tt.envVal)
			}
			s := slack.Slack{Channel: tt.channel}
			if got := s.GetWebHookURL(); got != tt.want {
				t.Errorf("GetWebHookURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCreateTaskMessage(t *testing.T) {
	due := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	cautionTasks := []task.Task{
		{Title: "タスクA", ShortURL: "https://short.url/a", Due: &due},
		{Title: "タスクB", ShortURL: "https://short.url/b", Due: nil},
	}
	deadTasks := []task.Task{
		{Title: "タスクC", ShortURL: "https://short.url/c", Due: &due},
	}

	s := slack.Slack{Channel: "00_today_tasks"}
	msg := s.CreateTaskMessage(cautionTasks, deadTasks)

	checks := []struct {
		desc string
		sub  string
	}{
		{"Key Tasks 見出し", "Key Tasks"},
		{"Dead Tasks 見出し", "Dead Tasks"},
		{"cautionTask Dueあり", "2024-06-15"},
		{"cautionTask Dueなし", "なるはや"},
		{"連番1", "1:"},
		{"連番2", "2:"},
		{"deadTask title", "タスクC"},
	}
	for _, c := range checks {
		if !strings.Contains(msg, c.sub) {
			t.Errorf("CreateTaskMessage() missing %s (%q)", c.desc, c.sub)
		}
	}
}

package task_test

import (
	"testing"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

func jst() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

func TestGetJSTDue(t *testing.T) {
	utc := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	tsk := task.Task{Due: &utc}
	got := tsk.GetJSTDue(&utc)
	wantHour := 9
	if got.Hour() != wantHour {
		t.Errorf("GetJSTDue hour = %d, want %d", got.Hour(), wantHour)
	}
	if got.Location().String() != "Asia/Tokyo" {
		t.Errorf("GetJSTDue location = %q, want \"Asia/Tokyo\"", got.Location().String())
	}
}

func TestIsDueOver(t *testing.T) {
	jst := jst()
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	past := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)
	today := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)
	future := time.Date(2024, 6, 16, 12, 0, 0, 0, jst)
	todayStart := time.Date(2024, 6, 15, 0, 0, 0, 0, jst)

	tests := []struct {
		name string
		due  *time.Time
		want bool
	}{
		{"過去", &past, true},
		{"今日の途中", &today, false},
		{"未来", &future, false},
		{"今日の開始（equal）", &todayStart, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task.Task{Due: tt.due}
			if got := tsk.IsDueOver(now); got != tt.want {
				t.Errorf("IsDueOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTodayTask(t *testing.T) {
	jst := jst()
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	todayMid := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)
	yesterday := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)
	tomorrow := time.Date(2024, 6, 16, 12, 0, 0, 0, jst)

	tests := []struct {
		name string
		due  *time.Time
		want bool
	}{
		{"今日のタスク", &todayMid, true},
		{"昨日のタスク", &yesterday, false},
		{"明日のタスク", &tomorrow, false},
		{"Due nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task.Task{Due: tt.due}
			if got := tsk.IsTodayTask(now); got != tt.want {
				t.Errorf("IsTodayTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

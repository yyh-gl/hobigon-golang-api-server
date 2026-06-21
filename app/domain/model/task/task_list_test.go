package task_test

import (
	"testing"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

func TestList_GetTodayTasks(t *testing.T) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	todayDue := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)
	yesterdayDue := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)
	tomorrowDue := time.Date(2024, 6, 16, 12, 0, 0, 0, jst)

	list := task.List{
		{Title: "今日のタスク", Due: &todayDue},
		{Title: "昨日のタスク", Due: &yesterdayDue},
		{Title: "明日のタスク", Due: &tomorrowDue},
		{Title: "Dueなし", Due: nil},
	}

	got := list.GetTodayTasks(now)
	if len(got) != 1 {
		t.Fatalf("GetTodayTasks() len = %d, want 1", len(got))
	}
	if got[0].Title != "今日のタスク" {
		t.Errorf("GetTodayTasks()[0].Title = %q, want \"今日のタスク\"", got[0].Title)
	}
}

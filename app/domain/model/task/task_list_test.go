// =============================================================================
// テストリスト（Canon TDD Step 1）
// 対象: List メソッド（GetToDoTasks, GetDoingTasks, GetDeadlineApproachingTasks, GetDueOverTasks）
//
// [List.GetToDoTasks()] — ステータス別取得
//   正常系:
//   - StatusToDo のみのリストから全件返す
//   - StatusToDo と StatusDoing 混在リストから StatusToDo のみ返す
//   境界値/特殊ケース:
//   - 空リストを渡すと空（len=0）を返す
//   - 一致するステータスがないリストを渡すと空（len=0）を返す
//
// [List.GetDoingTasks()] — ステータス別取得
//   正常系:
//   - StatusDoing のみのリストから全件返す
//   - StatusToDo と StatusDoing 混在リストから StatusDoing のみ返す
//   境界値/特殊ケース:
//   - 空リストを渡すと空（len=0）を返す
//   - 一致するステータスがないリストを渡すと空（len=0）を返す
//
// [List.GetDeadlineApproachingTasks(now)] — 新設
//   正常系:
//   - 期限間近タスクのみのリストから全件返す
//   - 期限切れタスクと期限間近タスクの混在リストから期限間近のみ返す
//   - Due が nil のタスクは含まない
//   境界値/特殊ケース:
//   - 空リストを渡すと空を返す
//   - 期限切れタスクのみのリストを渡すと空を返す
//   - 8 日後以降のタスクのみのリストを渡すと空を返す
//   - Due nil のタスクのみのリストを渡すと空を返す
//   - 今日 Due・7日後 Due・8日後 Due・期限切れ・Due nil の混在リストから
//     今日 Due と 7日後 Due のみ返す（境界値を含む両端）
//
// [List.GetDueOverTasks(now)] — 新設
//   正常系:
//   - 期限切れタスクのみのリストから全件返す
//   - 期限切れと期限間近の混在リストから期限切れのみ返す
//   - Due が nil のタスクは含まない
//   境界値/特殊ケース:
//   - 空リストを渡すと空を返す
//   - 期限間近タスクのみのリストを渡すと空を返す
//   - Due nil のタスクのみのリストを渡すと空を返す
//
// [排他性]
//   - 同一タスクが GetDeadlineApproachingTasks と GetDueOverTasks の両方に含まれないこと
//     （昨日 Due タスク: GetDueOverTasks に出現, GetDeadlineApproachingTasks に不出現）
//     （今日 Due タスク: GetDeadlineApproachingTasks に出現, GetDueOverTasks に不出現）
//
// =============================================================================

package task_test

import (
	"testing"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

func TestList_GetToDoTasks(t *testing.T) {
	todoOnly := task.List{
		{Title: "todo1", Status: task.StatusToDo},
		{Title: "todo2", Status: task.StatusToDo},
	}
	doingOnly := task.List{
		{Title: "doing1", Status: task.StatusDoing},
	}
	mixed := task.List{
		{Title: "todo", Status: task.StatusToDo},
		{Title: "doing", Status: task.StatusDoing},
	}

	t.Run("ToDoのみリストで全件返す", func(t *testing.T) {
		got := todoOnly.GetToDoTasks()
		if len(got) != 2 {
			t.Errorf("GetToDoTasks() len = %d, want 2", len(got))
		}
	})
	t.Run("混在リストからToDoのみ返す", func(t *testing.T) {
		got := mixed.GetToDoTasks()
		if len(got) != 1 || got[0].Title != "todo" {
			t.Errorf("GetToDoTasks() = %v, want [{todo}]", got)
		}
	})
	t.Run("空リストで空を返す", func(t *testing.T) {
		got := task.List{}.GetToDoTasks()
		if len(got) != 0 {
			t.Errorf("GetToDoTasks(empty) len = %d, want 0", len(got))
		}
	})
	t.Run("一致なしで空を返す", func(t *testing.T) {
		got := doingOnly.GetToDoTasks()
		if len(got) != 0 {
			t.Errorf("GetToDoTasks(no match) len = %d, want 0", len(got))
		}
	})
}

func TestList_GetDoingTasks(t *testing.T) {
	todoOnly := task.List{
		{Title: "todo1", Status: task.StatusToDo},
	}
	doingOnly := task.List{
		{Title: "doing1", Status: task.StatusDoing},
	}
	mixed := task.List{
		{Title: "todo", Status: task.StatusToDo},
		{Title: "doing", Status: task.StatusDoing},
	}

	t.Run("Doingのみリストで全件返す", func(t *testing.T) {
		got := doingOnly.GetDoingTasks()
		if len(got) != 1 {
			t.Errorf("GetDoingTasks() len = %d, want 1", len(got))
		}
	})
	t.Run("混在リストからDoingのみ返す", func(t *testing.T) {
		got := mixed.GetDoingTasks()
		if len(got) != 1 || got[0].Title != "doing" {
			t.Errorf("GetDoingTasks() = %v, want [{doing}]", got)
		}
	})
	t.Run("空リストで空を返す", func(t *testing.T) {
		got := task.List{}.GetDoingTasks()
		if len(got) != 0 {
			t.Errorf("GetDoingTasks(empty) len = %d, want 0", len(got))
		}
	})
	t.Run("一致なしで空を返す", func(t *testing.T) {
		got := todoOnly.GetDoingTasks()
		if len(got) != 0 {
			t.Errorf("GetDoingTasks(no match) len = %d, want 0", len(got))
		}
	})
}

func TestList_GetDeadlineApproachingTasks(t *testing.T) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	todayDue := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)
	in7Days := time.Date(2024, 6, 22, 12, 0, 0, 0, jst)
	in8Days := time.Date(2024, 6, 23, 12, 0, 0, 0, jst)
	yesterdayDue := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)

	t.Run("期限間近のみのリストから全件返す", func(t *testing.T) {
		list := task.List{
			{Title: "今日", Due: &todayDue},
			{Title: "7日後", Due: &in7Days},
		}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 2 {
			t.Errorf("GetDeadlineApproachingTasks() len = %d, want 2", len(got))
		}
	})
	t.Run("期限切れと期限間近の混在から期限間近のみ返す", func(t *testing.T) {
		list := task.List{
			{Title: "昨日", Due: &yesterdayDue},
			{Title: "今日", Due: &todayDue},
		}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 1 || got[0].Title != "今日" {
			t.Errorf("GetDeadlineApproachingTasks() = %v, want [今日]", got)
		}
	})
	t.Run("Due nilのタスクは含まない", func(t *testing.T) {
		list := task.List{
			{Title: "nilDue", Due: nil},
			{Title: "今日", Due: &todayDue},
		}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 1 || got[0].Title != "今日" {
			t.Errorf("GetDeadlineApproachingTasks() = %v, want [今日]", got)
		}
	})
	t.Run("空リストで空を返す", func(t *testing.T) {
		got := task.List{}.GetDeadlineApproachingTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDeadlineApproachingTasks(empty) len = %d, want 0", len(got))
		}
	})
	t.Run("期限切れのみで空を返す", func(t *testing.T) {
		list := task.List{{Title: "昨日", Due: &yesterdayDue}}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDeadlineApproachingTasks(overdue only) len = %d, want 0", len(got))
		}
	})
	t.Run("8日後以降のみで空を返す", func(t *testing.T) {
		list := task.List{{Title: "8日後", Due: &in8Days}}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDeadlineApproachingTasks(far future) len = %d, want 0", len(got))
		}
	})
	t.Run("Due nilのみで空を返す", func(t *testing.T) {
		list := task.List{{Title: "nilDue", Due: nil}}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDeadlineApproachingTasks(nil only) len = %d, want 0", len(got))
		}
	})
	t.Run("混在リストから今日と7日後のみ返す", func(t *testing.T) {
		list := task.List{
			{Title: "今日", Due: &todayDue},
			{Title: "7日後", Due: &in7Days},
			{Title: "8日後", Due: &in8Days},
			{Title: "昨日", Due: &yesterdayDue},
			{Title: "nilDue", Due: nil},
		}
		got := list.GetDeadlineApproachingTasks(now)
		if len(got) != 2 {
			t.Errorf("GetDeadlineApproachingTasks(mixed) len = %d, want 2", len(got))
		}
		if got[0].Title != "今日" || got[1].Title != "7日後" {
			t.Errorf("GetDeadlineApproachingTasks(mixed) titles = [%q, %q], want [\"今日\", \"7日後\"]", got[0].Title, got[1].Title)
		}
	})
}

func TestList_GetDueOverTasks(t *testing.T) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	todayDue := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)
	yesterdayDue := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)
	twoDaysAgo := time.Date(2024, 6, 13, 12, 0, 0, 0, jst)

	t.Run("期限切れのみのリストから全件返す", func(t *testing.T) {
		list := task.List{
			{Title: "昨日", Due: &yesterdayDue},
			{Title: "2日前", Due: &twoDaysAgo},
		}
		got := list.GetDueOverTasks(now)
		if len(got) != 2 {
			t.Errorf("GetDueOverTasks() len = %d, want 2", len(got))
		}
	})
	t.Run("期限切れと期限間近の混在から期限切れのみ返す", func(t *testing.T) {
		list := task.List{
			{Title: "昨日", Due: &yesterdayDue},
			{Title: "今日", Due: &todayDue},
		}
		got := list.GetDueOverTasks(now)
		if len(got) != 1 || got[0].Title != "昨日" {
			t.Errorf("GetDueOverTasks() = %v, want [昨日]", got)
		}
	})
	t.Run("Due nilのタスクは含まない", func(t *testing.T) {
		list := task.List{
			{Title: "nilDue", Due: nil},
			{Title: "昨日", Due: &yesterdayDue},
		}
		got := list.GetDueOverTasks(now)
		if len(got) != 1 || got[0].Title != "昨日" {
			t.Errorf("GetDueOverTasks() = %v, want [昨日]", got)
		}
	})
	t.Run("空リストで空を返す", func(t *testing.T) {
		got := task.List{}.GetDueOverTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDueOverTasks(empty) len = %d, want 0", len(got))
		}
	})
	t.Run("期限間近のみで空を返す", func(t *testing.T) {
		list := task.List{{Title: "今日", Due: &todayDue}}
		got := list.GetDueOverTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDueOverTasks(approaching only) len = %d, want 0", len(got))
		}
	})
	t.Run("Due nilのみで空を返す", func(t *testing.T) {
		list := task.List{{Title: "nilDue", Due: nil}}
		got := list.GetDueOverTasks(now)
		if len(got) != 0 {
			t.Errorf("GetDueOverTasks(nil only) len = %d, want 0", len(got))
		}
	})
}

func TestList_GetDeadlineApproachingAndGetDueOverExclusive(t *testing.T) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	yesterdayDue := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)
	todayDue := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)

	list := task.List{
		{Title: "昨日", Due: &yesterdayDue},
		{Title: "今日", Due: &todayDue},
	}

	t.Run("昨日DueはGetDueOverTasksに出現しGetDeadlineApproachingTasksに不出現", func(t *testing.T) {
		dueOver := list.GetDueOverTasks(now)
		approaching := list.GetDeadlineApproachingTasks(now)

		foundInDueOver := false
		for _, t := range dueOver {
			if t.Title == "昨日" {
				foundInDueOver = true
			}
		}
		foundInApproaching := false
		for _, t := range approaching {
			if t.Title == "昨日" {
				foundInApproaching = true
			}
		}
		if !foundInDueOver {
			t.Error("昨日DueはGetDueOverTasksに含まれるべき")
		}
		if foundInApproaching {
			t.Error("昨日DueはGetDeadlineApproachingTasksに含まれてはいけない")
		}
	})
	t.Run("今日DueはGetDeadlineApproachingTasksに出現しGetDueOverTasksに不出現", func(t *testing.T) {
		dueOver := list.GetDueOverTasks(now)
		approaching := list.GetDeadlineApproachingTasks(now)

		foundInDueOver := false
		for _, t := range dueOver {
			if t.Title == "今日" {
				foundInDueOver = true
			}
		}
		foundInApproaching := false
		for _, t := range approaching {
			if t.Title == "今日" {
				foundInApproaching = true
			}
		}
		if foundInDueOver {
			t.Error("今日DueはGetDueOverTasksに含まれてはいけない")
		}
		if !foundInApproaching {
			t.Error("今日DueはGetDeadlineApproachingTasksに含まれるべき")
		}
	})
}

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

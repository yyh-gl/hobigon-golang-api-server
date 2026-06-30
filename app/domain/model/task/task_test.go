// =============================================================================
// テストリスト（Canon TDD Step 1）
// 対象: Task メソッド（IsDueOver, IsDeadlineApproaching）
//
// [Task.IsDueOver(now)] — nil セーフ修正後の挙動
//   - Due が nil のとき false を返す（panic しない）
//   ※ 以下は既存テスト（TestIsDueOver）でカバー済み
//   - Due が昨日の日時のとき true を返す
//   - Due が今日の途中の日時のとき false を返す（今日は期限切れでない）
//   - Due が今日の 0 時（todayStart）のとき false を返す（equal は期限切れでない）
//   - Due が明日以降の日時のとき false を返す
//
// [Task.IsDeadlineApproaching(now)] — 新設（今日 ≤ Due ≤ 今日+7日）
//   正常系:
//   - Due が今日の日付内のとき true を返す
//   - Due が 7 日後の日付内のとき true を返す
//   境界値:
//   - Due が今日の 0 時（todayStart）のとき true を返す（下限を含む）
//   - Due が今日の 23:59:59 のとき true を返す（今日の末尾は含む）
//   - Due が 7 日後の 23:59:59 のとき true を返す（上限末尾は含む）
//   - Due が 8 日後の 0 時のとき false を返す（上限の外、8日後の0時は含まない）
//   異常系/特殊ケース:
//   - Due が nil のとき false を返す（nil セーフ）
//   - Due が昨日（過去）のとき false を返す（期限切れは期限間近でない）
//   - Due が 8 日後以降のとき false を返す
//   排他性:
//   - 昨日 Due のタスクは IsDueOver=true かつ IsDeadlineApproaching=false（同時に true にならない）
//   - 今日 Due のタスクは IsDueOver=false かつ IsDeadlineApproaching=true
//
// =============================================================================

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
	tsk := task.Task{Deadline: &utc}
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
		{"Due nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task.Task{Deadline: tt.due}
			if got := tsk.IsDead(now); got != tt.want {
				t.Errorf("IsDueOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDeadlineApproaching(t *testing.T) {
	jst := jst()
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)

	todayStart := time.Date(2024, 6, 15, 0, 0, 0, 0, jst)
	todayMid := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)
	todayEnd := time.Date(2024, 6, 15, 23, 59, 59, 0, jst)
	in7Days := time.Date(2024, 6, 22, 12, 0, 0, 0, jst)
	in7DaysEnd := time.Date(2024, 6, 22, 23, 59, 59, 0, jst)
	in8DaysStart := time.Date(2024, 6, 23, 0, 0, 0, 0, jst)
	in8Days := time.Date(2024, 6, 23, 12, 0, 0, 0, jst)
	yesterday := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)

	tests := []struct {
		name string
		due  *time.Time
		want bool
	}{
		{"今日の日付内", &todayMid, true},
		{"7日後の日付内", &in7Days, true},
		{"今日の0時（下限）", &todayStart, true},
		{"今日の23:59:59", &todayEnd, true},
		{"7日後の23:59:59（上限末尾）", &in7DaysEnd, true},
		{"8日後の0時（上限の外）", &in8DaysStart, false},
		{"Due nil", nil, false},
		{"昨日（過去）", &yesterday, false},
		{"8日後以降", &in8Days, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task.Task{Deadline: tt.due}
			if got := tsk.IsDeadlineApproaching(now); got != tt.want {
				t.Errorf("IsDeadlineApproaching() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDeadlineApproachingAndIsDueOverExclusive(t *testing.T) {
	jst := jst()
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, jst)
	yesterday := time.Date(2024, 6, 14, 12, 0, 0, 0, jst)
	today := time.Date(2024, 6, 15, 12, 0, 0, 0, jst)

	t.Run("昨日DueはIsDueOver=trueかつIsDeadlineApproaching=false", func(t *testing.T) {
		tsk := task.Task{Deadline: &yesterday}
		if !tsk.IsDead(now) {
			t.Error("IsDueOver() = false, want true")
		}
		if tsk.IsDeadlineApproaching(now) {
			t.Error("IsDeadlineApproaching() = true, want false")
		}
	})
	t.Run("今日DueはIsDueOver=falseかつIsDeadlineApproaching=true", func(t *testing.T) {
		tsk := task.Task{Deadline: &today}
		if tsk.IsDead(now) {
			t.Error("IsDueOver() = true, want false")
		}
		if !tsk.IsDeadlineApproaching(now) {
			t.Error("IsDeadlineApproaching() = false, want true")
		}
	})
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
			tsk := task.Task{Deadline: tt.due}
			if got := tsk.IsTodayTask(now); got != tt.want {
				t.Errorf("IsTodayTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

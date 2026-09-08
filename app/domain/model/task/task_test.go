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
//
// テストリスト（Canon TDD Step 1: TZ正規化）
// 対象: Task メソッド（IsDead, IsDeadlineApproaching, IsTodayTask）の`now`TZ正規化
//
// [正常系・回帰]（`now`がUTC＝本番相当のとき）
//   - now=2024-06-14 22:00 UTC（06-15 07:00JST）、Deadline=2024-06-14 00:00 UTC（昨日・日付のみ）
//     → IsDead=true, IsDeadlineApproaching=false
//   - 同now、Deadline=2024-06-15 00:00 UTC（今日）→ IsDead=false, IsDeadlineApproaching=true
//   - 同now、Deadline=2024-06-22 00:00 UTC（7日後）→ IsDeadlineApproaching=true
//   - 同now、Deadline=2024-06-23 00:00 UTC（8日後）→ IsDeadlineApproaching=false
//   - 同now、Deadline=2024-06-14T23:00:00+09:00（時刻付き・JST基準で昨日23時）→ IsDead=true
//   - 同now、Deadline=2024-06-15T01:00:00+09:00（時刻付き・JST基準で今日1時）
//     → IsDead=false, IsDeadlineApproaching=true
//   - now がUTCのとき、IsTodayTask が今日のDeadlineに対してtrueを返す
//   - now がUTCのとき、IsTodayTask が昨日のDeadlineに対してfalseを返す
//   - now がUTCのとき、IsTodayTask が明日のDeadlineに対してfalseを返す
//
// [境界値]（TZ正規化ヘルパーの本当の境界＝JST 0時）
//   - now=2024-06-14 15:00:00 UTC（06-15 00:00 JSTちょうど）、Deadline=2024-06-14 00:00 UTC
//     → IsDead=true
//   - now=2024-06-14 14:59:59 UTC（06-14 23:59:59 JST・上記の1秒前）、同Deadline
//     → 「今日」扱いになりIsDead=false, IsDeadlineApproaching=true
//
// [特殊ケース]（TZ非依存性）
//   - 同一瞬間のnowをJST/UTC/任意の別TZ（FixedZone("X", -5*60*60)）で与えても、
//     同じDeadlineに対するIsDead/IsDeadlineApproaching/IsTodayTaskの結果が一致する
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

func TestIsDead_NowIsUTC(t *testing.T) {
	now := time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC)

	tests := []struct {
		name                      string
		deadline                  time.Time
		wantIsDead                bool
		wantIsDeadlineApproaching bool
	}{
		{"昨日Deadline（日付のみ）", time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC), true, false},
		{"今日Deadline（日付のみ）", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false, true},
		{"7日後Deadline（日付のみ）", time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC), false, true},
		{"8日後Deadline（日付のみ）", time.Date(2024, 6, 23, 0, 0, 0, 0, time.UTC), false, false},
		{"昨日23時JST（時刻付き）", parseJST(t, "2024-06-14T23:00:00+09:00"), true, false},
		{"今日1時JST（時刻付き）", parseJST(t, "2024-06-15T01:00:00+09:00"), false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tsk := task.Task{Deadline: &tt.deadline}
			if got := tsk.IsDead(now); got != tt.wantIsDead {
				t.Errorf("IsDead() = %v, want %v", got, tt.wantIsDead)
			}
			if got := tsk.IsDeadlineApproaching(now); got != tt.wantIsDeadlineApproaching {
				t.Errorf("IsDeadlineApproaching() = %v, want %v", got, tt.wantIsDeadlineApproaching)
			}
		})
	}
}

func TestIsTodayTask_NowIsUTC(t *testing.T) {
	now := time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC)

	today := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	yesterday := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)
	tomorrow := time.Date(2024, 6, 16, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		due  *time.Time
		want bool
	}{
		{"今日のタスク", &today, true},
		{"昨日のタスク", &yesterday, false},
		{"明日のタスク", &tomorrow, false},
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

func TestIsDead_JSTMidnightBoundary(t *testing.T) {
	deadline := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)

	t.Run("JST0時ちょうど（昨日Deadlineが期限切れになる）", func(t *testing.T) {
		now := time.Date(2024, 6, 14, 15, 0, 0, 0, time.UTC)
		tsk := task.Task{Deadline: &deadline}
		if !tsk.IsDead(now) {
			t.Error("IsDead() = false, want true")
		}
	})
	t.Run("JST0時の1秒前（まだ今日扱い）", func(t *testing.T) {
		now := time.Date(2024, 6, 14, 14, 59, 59, 0, time.UTC)
		tsk := task.Task{Deadline: &deadline}
		if tsk.IsDead(now) {
			t.Error("IsDead() = true, want false")
		}
		if !tsk.IsDeadlineApproaching(now) {
			t.Error("IsDeadlineApproaching() = false, want true")
		}
	})
}

func TestIsDead_IndependentOfNowLocation(t *testing.T) {
	deadline := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)
	nowUTC := time.Date(2024, 6, 14, 22, 0, 0, 0, time.UTC)
	other := time.FixedZone("X", -5*60*60)

	nows := []time.Time{
		nowUTC,
		nowUTC.In(jst()),
		nowUTC.In(other),
	}

	tsk := task.Task{Deadline: &deadline}
	for i, now := range nows {
		if i == 0 {
			continue
		}
		if got, want := tsk.IsDead(now), tsk.IsDead(nows[0]); got != want {
			t.Errorf("IsDead() with now[%d] = %v, want %v (same as UTC)", i, got, want)
		}
		if got, want := tsk.IsDeadlineApproaching(now), tsk.IsDeadlineApproaching(nows[0]); got != want {
			t.Errorf("IsDeadlineApproaching() with now[%d] = %v, want %v (same as UTC)", i, got, want)
		}
		if got, want := tsk.IsTodayTask(now), tsk.IsTodayTask(nows[0]); got != want {
			t.Errorf("IsTodayTask() with now[%d] = %v, want %v (same as UTC)", i, got, want)
		}
	}
}

func parseJST(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("parseJST(%q): %v", s, err)
	}
	return tm
}

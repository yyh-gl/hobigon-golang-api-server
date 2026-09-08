package task

import (
	"time"
)

type Status string

const (
	StatusToDo  Status = "To Do"
	StatusDoing Status = "Doing"
)

func (s Status) String() string {
	return string(s)
}

// Task : タスクを表すドメインモデル
// TODO: ドメインモデル貧血症を治す
type Task struct {
	ID            string     `json:"-"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Deadline      *time.Time `json:"deadline"`
	Status        Status     `json:"status"`
	List          string     `json:"list"`
	ShortURL      string     `json:"short_url"`
	OriginalModel any        `json:"-"`
}

// GetJSTDue : 日本時間の期限を取得
func (t Task) GetJSTDue(utcDue *time.Time) *time.Time {
	jst := getJSTNow()
	jstDue := utcDue.In(jst)
	return &jstDue
}

// IsDead : 期限切れかどうか判定
func (t Task) IsDead(now time.Time) (isDueOver bool) {
	if t.Deadline == nil {
		return false
	}
	todayStart := todayStartJST(now)
	return !t.Deadline.Equal(todayStart) && t.Deadline.Before(todayStart)
}

// IsDeadlineApproaching : 期限が近づいているかどうか判定（今日 ≤ Deadline ≤ 今日+7日）
func (t Task) IsDeadlineApproaching(now time.Time) bool {
	if t.Deadline == nil {
		return false
	}
	todayStart := todayStartJST(now)
	return !t.Deadline.Before(todayStart) && t.Deadline.Before(todayStart.AddDate(0, 0, 8))
}

// IsTodayTask : 今日のタスクかどうか判定
func (t Task) IsTodayTask(now time.Time) (isTodayTask bool) {
	todayStart := todayStartJST(now)
	todayEnd := todayStart.AddDate(0, 0, 1).Add(-time.Second)
	if t.Deadline != nil && t.Deadline.After(todayStart) && t.Deadline.Before(todayEnd) {
		return true
	}
	return false
}

// getJSTNow : 現在時刻を日本時間で取得
func getJSTNow() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

// todayStartJST : `now`をJSTに正規化した上での「今日の0時JST」を算出
func todayStartJST(now time.Time) time.Time {
	jst := getJSTNow()
	n := now.In(jst)
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, jst)
}

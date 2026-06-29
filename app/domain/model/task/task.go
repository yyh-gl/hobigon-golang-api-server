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
	Due           *time.Time `json:"due"`
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

// IsDueOver : 期限切れかどうか判定
func (t Task) IsDueOver(now time.Time) (isDueOver bool) {
	if t.Due == nil {
		return false
	}
	jst := getJSTNow()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	return !t.Due.Equal(todayStart) && t.Due.Before(todayStart)
}

// IsDeadlineApproaching : 期限が近づいているかどうか判定（今日 ≤ Due ≤ 今日+7日）
func (t Task) IsDeadlineApproaching(now time.Time) bool {
	if t.Due == nil {
		return false
	}
	jst := getJSTNow()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	return !t.Due.Before(todayStart) && t.Due.Before(todayStart.AddDate(0, 0, 8))
}

// IsTodayTask : 今日のタスクかどうか判定
func (t Task) IsTodayTask(now time.Time) (isTodayTask bool) {
	jst := getJSTNow()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, jst)
	if t.Due != nil && t.Due.After(todayStart) && t.Due.Before(todayEnd) {
		return true
	}
	return false
}

// getJSTNow : 現在時刻を日本時間で取得
func getJSTNow() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

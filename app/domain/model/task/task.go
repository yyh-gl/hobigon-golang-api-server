package task

import (
	"time"
)

// Task : タスク用のドメインモデル
type Task struct {
	Title         string      `json:"title"`
	Description   string      `json:"description"`
	Due           *time.Time  `json:"due"`
	Board         string      `json:"board"`
	List          string      `json:"list"`
	ShortURL      string      `json:"short_url"`
	OriginalModel interface{} `json:"-"`
}

// GetJSTDue : 日本時間の期限を取得
func (t Task) GetJSTDue(utcDue *time.Time) *time.Time {
	jst := getJSTNow()
	jstDue := utcDue.In(jst)
	return &jstDue
}

// IsDueOver : 期限切れかどうか判定
func (t Task) IsDueOver() (isDueOver bool) {
	jst := getJSTNow()
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	return !t.Due.Equal(todayStart) && t.Due.Before(todayStart)
}

// IsTodayTask : 今日のタスクかどうか判定
func (t Task) IsTodayTask() (isTodayTask bool) {
	jst := getJSTNow()
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	todayEnd := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, jst)
	if t.Due != nil && t.Due.After(todayStart) && t.Due.Before(todayEnd) {
		return true
	}
	return false
}

// getJSTNow : 現在時刻を日本時間で取得
func getJSTNow() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

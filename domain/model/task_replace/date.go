package task

import "time"

// Date : Amazon の欲しい物リストを表す値オブジェクト
type Date struct {
	value *time.Time
}

// NewDate : Date を生成
func NewDate(val *time.Time) (*Date, error) {
	jst := getJST()
	jstDate := val.In(jst)
	return &Date{value: &jstDate}, nil
}

// IsDueOver : 期限切れかどうか判定
func (d Date) IsDueOver() bool {
	jst := getJST()
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	return !d.value.Equal(todayStart) && d.value.Before(todayStart)
}

// IsTodayTask : 今日のタスクかどうか判定
func (d Date) IsTodayTask() bool {
	jst := getJST()
	today := time.Now()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, jst)
	todayEnd := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, jst)
	if d.value != nil && d.value.After(todayStart) && d.value.Before(todayEnd) {
		return true
	}
	return false
}

// getJSTNow : 日本標準時のロケーション情報を取得
func getJST() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}

package birthday

import (
	"time"
)

// time.Time のレイアウト
const layout = "0102"

// Date : Amazon の欲しい物リストを表す値オブジェクト
type Date struct {
	value string
}

// NewDate : Date を生成
func NewDate(val time.Time) (*Date, error) {
	return &Date{value: val.Format(layout)}, nil
}

// String : Date の値を文字列として返却
func (d Date) String() string {
	return d.value
}

// IsToday : 指定 Date が本日かどうか判定
func (d Date) IsToday() bool {
	today := time.Now().Format(layout)
	return d.value == today
}

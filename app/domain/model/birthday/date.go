package birthday

import (
	"errors"
	"time"
)

const (
	// 日付（文字列）の長さ
	lengthDate = 4
	// time.Time→文字列のレイアウト
	layout = "0102"
)

// Date : 誕生日に関する日付を表す値オブジェクト
type Date string

// NewDate : Date を生成
func NewDate(val string) (*Date, error) {
	if len(val) != lengthDate {
		return nil, errors.New("バリデーションエラー：【Birthday】Dateの形式が誤っています")
	}

	d := Date(val)
	return &d, nil
}

// String : Date の値を文字列として返却
func (d Date) String() string {
	return string(d)
}

// Equal : 同値判定
func (d Date) Equal(date Date) bool {
	return d == date
}

// IsToday : 指定 Date が本日かどうか判定
func (d Date) IsToday() bool {
	today := time.Now().Format(layout)
	return d.Equal(Date(today))
}

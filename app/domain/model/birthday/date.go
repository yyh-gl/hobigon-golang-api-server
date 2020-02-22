package birthday

import (
	"errors"
	"time"
)

const (
	// defaultLength : 日付（文字列）の長さ
	defaultLength = 4
	// defaultLayout : time.Time→文字列のレイアウト
	defaultLayout = "0102"
)

// Date : 誕生日に関する日付を表す値オブジェクト
type Date string

// newDate : Date を生成
func newDate(val string) (*Date, error) {
	if len(val) != defaultLength {
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
	today := time.Now().Format(defaultLayout)
	return d.Equal(Date(today))
}

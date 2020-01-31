package blog

import (
	"errors"
)

// maxLength : 名前の長さ
const maxLength = 30

// Title : ブログのタイトルを表す値オブジェクト
type Title string

// NewTitle : Title を生成
func NewTitle(val string) (*Title, error) {
	if len(val) > maxLength {
		return nil, errors.New("バリデーションエラー：【Blog】Titleは30文字以内です")
	}

	n := Title(val)
	return &n, nil
}

// String : Title の値を文字列として返却
func (n Title) String() string {
	return string(n)
}

// IsNull : Title の値が Null かどうか判定
func (n Title) IsNull() bool {
	return n == ""
}

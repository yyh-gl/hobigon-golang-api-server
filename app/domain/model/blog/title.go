package blog

import (
	"errors"
	"strconv"
)

// maxLength : 名前の長さ
const maxLength = 50

// Title : ブログのタイトルを表す値オブジェクト
type Title string

// newTitle : Titleを生成
func newTitle(val string) (Title, error) {
	if len(val) > maxLength {
		return "", errors.New("バリデーションエラー：【Blog】Titleは" + strconv.Itoa(maxLength) + "文字以内です")
	}

	return Title(val), nil
}

// String : Titleの値を文字列として返却
func (n Title) String() string {
	return string(n)
}

// IsNull : Titleの値がNullかどうか判定
func (n Title) IsNull() bool {
	return n == ""
}

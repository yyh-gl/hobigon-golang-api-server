package birthday

import (
	"errors"
)

// maxLength : 名前の長さ
const maxLength = 30

// Name : 誕生日の人の名前を表す値オブジェクト
type Name string

// NewName : Name を生成
func NewName(val string) (*Name, error) {
	if len(val) > maxLength {
		return nil, errors.New("バリデーションエラー：【Birthday】Nameは30文字以内です")
	}

	n := Name(val)
	return &n, nil
}

// String : Name の値を文字列として返却
func (n Name) String() string {
	return string(n)
}

// IsNull : Name の値が Null かどうか判定
func (n Name) IsNull() bool {
	return n == ""
}

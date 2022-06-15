package birthday

import (
	"errors"
	"strconv"
)

// maxLength : 名前の長さ
const maxLength = 30

// Name : 誕生日の人の名前を表す値オブジェクト
type Name string

// newName : Name を生成
func newName(val string) (Name, error) {
	if len(val) > maxLength {
		return "", errors.New("バリデーションエラー：【Birthday】Nameは" + strconv.Itoa(maxLength) + "文字以内です")
	}

	return Name(val), nil
}

// String : Name の値を文字列として返却
func (n Name) String() string {
	return string(n)
}

// IsNull : Name の値が Null かどうか判定
func (n Name) IsNull() bool {
	return n == ""
}

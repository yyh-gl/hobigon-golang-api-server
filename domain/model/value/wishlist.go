package value

import (
	"errors"
	"strings"
)

// WishList : Amazon の欲しい物リストを表す値オブジェクト
type WishList struct {
	value string
}

// NewWishList : WishList を生成
func NewWishList(val string) (*WishList, error) {
	// WishList が空じゃない場合は "https://" で始まっていることをチェック
	if val != "" && !strings.HasPrefix(val, "https://") {
		return nil, errors.New("バリデーションエラー：【Birthday】WishList が \"https://\" から始まっていません")
	}

	return &WishList{value: val}, nil
}

// String : WishList の値を文字列として返却
func (wl WishList) String() string {
	return wl.value
}

func (wl WishList) IsNull() bool {
	return wl.value == ""
}

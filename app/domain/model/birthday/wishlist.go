package birthday

import (
	"errors"
	"strings"
)

// WishList : Amazon の欲しい物リストを表す値オブジェクト
type WishList string

// NewWishList : WishList を生成
func NewWishList(val string) (*WishList, error) {
	// WishList が空じゃない場合は "https://" で始まっていることをチェック
	if val != "" && !strings.HasPrefix(val, "https://") {
		return nil, errors.New("バリデーションエラー：【Birthday】WishListが\"https://\"から始まっていません")
	}

	wl := WishList(val)
	return &wl, nil
}

// String : WishList の値を文字列として返却
func (wl WishList) String() string {
	return string(wl)
}

// IsNull : WishList の値が Null かどうか判定
func (wl WishList) IsNull() bool {
	return wl == ""
}

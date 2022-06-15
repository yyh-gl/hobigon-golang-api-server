package birthday

import (
	"errors"
	"strings"
)

// WishList : Amazonの欲しい物リストを表す値オブジェクト
type WishList string

// newWishList : WishListを生成
func newWishList(val string) (WishList, error) {
	// WishList が空じゃない場合は "https://" で始まっていることをチェック
	if val != "" && !strings.HasPrefix(val, "https://") {
		return "", errors.New("バリデーションエラー：【Birthday】WishListが\"https://\"から始まっていません")
	}

	return WishList(val), nil
}

// String : WishListの値を文字列として返却
func (wl WishList) String() string {
	return string(wl)
}

// IsNull : WishListの値がNullかどうか判定
func (wl WishList) IsNull() bool {
	return wl == ""
}

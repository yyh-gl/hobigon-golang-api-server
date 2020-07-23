package birthday

import (
	"fmt"
)

// Birthday : 誕生日を表すドメインモデル
type Birthday struct {
	name     Name
	date     Date
	wishList WishList
}

// NewBirthday : Birthdayドメインモデルを生成
func NewBirthday(name string, date string, wishList string) (*Birthday, error) {
	// Nameを生成
	n, err := newName(name)
	if err != nil {
		return nil, fmt.Errorf("NewName()内でエラー: %w", err)
	}

	// Dateを生成
	d, err := newDate(date)
	if err != nil {
		return nil, fmt.Errorf("NewDate()内でエラー: %w", err)
	}

	// WishListを生成
	wl, err := newWishList(wishList)
	if err != nil {
		return nil, fmt.Errorf("NewWishList()内でエラー: %w", err)
	}

	return &Birthday{
		name:     *n,
		date:     *d,
		wishList: *wl,
	}, nil
}

// Name : name のゲッター
func (b Birthday) Name() Name {
	return b.name
}

// Date : date のゲッター
func (b Birthday) Date() Date {
	return b.date
}

// WishList : wishList のゲッター
func (b Birthday) WishList() WishList {
	return b.wishList
}

// CreateBirthdayMessage : 誕生日メッセージを生成
func (b Birthday) CreateBirthdayMessage() string {
	wishList := b.wishList.String()
	if b.wishList.IsNull() {
		wishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.name.String() + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + wishList + " :gainings:"
}

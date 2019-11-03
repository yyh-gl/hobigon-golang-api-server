package model

import (
	"time"

	"github.com/pkg/errors"
)

type Birthday struct {
	id        uint
	name      string
	date      Date
	wishList  WishList
	createdAt *time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

// NewBirthday : Birthday ドメインモデルを生成
func NewBirthday(name string, date time.Time, wishList string) (*Birthday, error) {
	// Date を生成
	d, err := NewDate(date)
	if err != nil {
		return nil, errors.Wrap(err, "NewDate()内でのエラー")
	}

	// WishList を生成
	wl, err := NewWishList(wishList)
	if err != nil {
		return nil, errors.Wrap(err, "NewWishList()内でのエラー")
	}

	return &Birthday{
		name:     name,
		date:     *d,
		wishList: *wl,
	}, nil
}

// SetID : id のセッター
func (b *Birthday) SetID(id uint) {
	b.id = id
}

// ID : id のゲッター
func (b Birthday) ID() uint {
	return b.id
}

// SetName : name のセッター
func (b *Birthday) SetName(name string) {
	b.name = name
}

// Name : name のゲッター
func (b Birthday) Name() string {
	return b.name
}

// SetDate : date のセッター
func (b *Birthday) SetDate(date Date) {
	b.date = date
}

// Date : date のゲッター
func (b Birthday) Date() Date {
	return b.date
}

// SetWishList : wishList のセッター
func (b *Birthday) SetWishList(wishList WishList) {
	b.wishList = wishList
}

// WishList : wishList のゲッター
func (b Birthday) WishList() WishList {
	return b.wishList
}

// SetCreatedAt : createdAt のセッター
func (b *Birthday) SetCreatedAt(createdAt *time.Time) {
	b.createdAt = createdAt
}

// CreatedAt : createdAt のゲッター
func (b Birthday) CreatedAt() *time.Time {
	return b.createdAt
}

// SetUpdatedAt : updatedAt のセッター
func (b *Birthday) SetUpdatedAt(updatedAt *time.Time) {
	b.updatedAt = updatedAt
}

// UpdatedAt : updatedAt のゲッター
func (b Birthday) UpdatedAt() *time.Time {
	return b.updatedAt
}

// SetDeletedAt : deletedAt のセッター
func (b *Birthday) SetDeletedAt(deletedAt *time.Time) {
	b.deletedAt = deletedAt
}

// DeletedAt : deletedAt のゲッター
func (b Birthday) DeletedAt() *time.Time {
	return b.deletedAt
}

// CreateBirthdayMessage : 誕生日メッセージを生成
func (b Birthday) CreateBirthdayMessage() string {
	wishList := b.wishList.String()
	if b.wishList.IsNull() {
		wishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.name + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + wishList + " :gainings:"
}

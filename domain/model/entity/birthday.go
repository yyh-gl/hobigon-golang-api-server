package entity

import (
	"time"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/value"
)

type Birthday struct {
	id        uint
	name      string
	date      value.Date
	wishList  value.WishList
	createdAt *time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

// NewBirthday : Birthday ドメインモデルを生成
func NewBirthday(name string, date time.Time, wishList string) (*Birthday, error) {
	// Date を生成
	d, err := value.NewDate(date)
	if err != nil {
		return nil, errors.Wrap(err, "NewDate()内でのエラー")
	}

	// WishList を生成
	wl, err := value.NewWishList(wishList)
	if err != nil {
		return nil, errors.Wrap(err, "NewWishList()内でのエラー")
	}

	return &Birthday{
		name:     name,
		date:     *d,
		wishList: *wl,
	}, nil
}

// NewBirthdayWithFullParams : パラメータ全指定で Birthday ドメインモデルを生成
func NewBirthdayWithFullParams(
	id uint,
	name string,
	date time.Time,
	wishList string,
	createdAt *time.Time,
	updatedAt *time.Time,
	deletedAt *time.Time,
) (*Birthday, error) {
	// Date を生成
	d, err := value.NewDate(date)
	if err != nil {
		return nil, errors.Wrap(err, "NewDate()内でのエラー")
	}

	// WishList を生成
	wl, err := value.NewWishList(wishList)
	if err != nil {
		return nil, errors.Wrap(err, "NewWishList()内でのエラー")
	}

	return &Birthday{
		id:        id,
		name:      name,
		date:      *d,
		wishList:  *wl,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}, nil
}

// ID : id のゲッター
func (b Birthday) ID() uint {
	return b.id
}

// Name : name のゲッター
func (b Birthday) Name() string {
	return b.name
}

// Date : date のゲッター
func (b Birthday) Date() value.Date {
	return b.date
}

// WishList : wishList のゲッター
func (b Birthday) WishList() value.WishList {
	return b.wishList
}

// CreatedAt : createdAt のゲッター
func (b Birthday) CreatedAt() *time.Time {
	return b.createdAt
}

// UpdatedAt : updatedAt のゲッター
func (b Birthday) UpdatedAt() *time.Time {
	return b.updatedAt
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

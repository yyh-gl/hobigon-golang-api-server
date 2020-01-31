package birthday

import (
	"fmt"
	"time"
)

// TODO: ドメインモデル貧血症を治す

// Birthday : 誕生日を表すドメインモデル
type Birthday struct {
	fields
}

type fields struct {
	ID        uint       `json:"id,omitempty"`
	Name      Name       `json:"name,omitempty"`
	Date      Date       `json:"date,omitempty"`
	WishList  WishList   `json:"wish_list,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewBirthday : Birthdayドメインモデルを生成
func NewBirthday(name string, date string, wishList string) (*Birthday, error) {
	// Name を生成
	n, err := NewName(name)
	if err != nil {
		return nil, fmt.Errorf("NewName()内でエラー: %w", err)
	}

	// Date を生成
	d, err := NewDate(date)
	if err != nil {
		return nil, fmt.Errorf("NewDate()内でエラー: %w", err)
	}

	// WishList を生成
	wl, err := NewWishList(wishList)
	if err != nil {
		return nil, fmt.Errorf("NewWishList()内でエラー: %w", err)
	}

	return &Birthday{
		fields{
			Name:     *n,
			Date:     *d,
			WishList: *wl,
		},
	}, nil
}

// NewBirthdayWithFullParams : パラメータ全指定でBirthdayドメインモデルを生成
func NewBirthdayWithFullParams(
	id uint,
	name string,
	date string,
	wishList string,
	createdAt *time.Time,
	updatedAt *time.Time,
	deletedAt *time.Time,
) (*Birthday, error) {
	// Name を生成
	n, err := NewName(name)
	if err != nil {
		return nil, fmt.Errorf("NewName()内でエラー: %w", err)
	}

	// Date を生成
	d, err := NewDate(date)
	if err != nil {
		return nil, fmt.Errorf("NewDate()内でエラー: %w", err)
	}

	// WishList を生成
	wl, err := NewWishList(wishList)
	if err != nil {
		return nil, fmt.Errorf("NewWishList()内でエラー: %w", err)
	}

	return &Birthday{
		fields{
			ID:        id,
			Name:      *n,
			Date:      *d,
			WishList:  *wl,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		},
	}, nil
}

// Name : name のゲッター
func (b Birthday) Name() Name {
	return b.fields.Name
}

// Date : date のゲッター
func (b Birthday) Date() Date {
	return b.fields.Date
}

// WishList : wishList のゲッター
func (b Birthday) WishList() WishList {
	return b.fields.WishList
}

// CreateBirthdayMessage : 誕生日メッセージを生成
func (b Birthday) CreateBirthdayMessage() string {
	wishList := b.fields.WishList.String()
	if b.fields.WishList.IsNull() {
		wishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.fields.Name.String() + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + wishList + " :gainings:"
}

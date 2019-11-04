package entity

import (
	"time"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/value"
)

// Birthday : 誕生日用のドメインモデル
type Birthday struct {
	id        uint
	name      string
	date      value.Date
	wishList  value.WishList
	createdAt *time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

type birthdayJSONFields struct {
	ID        uint       `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Date      string     `json:"date,omitempty"`
	WishList  string     `json:"wish_list,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// BirthdayJSON : 誕生日用の JSON レスポンス形式の定義
type BirthdayJSON struct {
	birthdayJSONFields
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

// JSONSerialize : JSON タグを含む構造体を返却
func (b Birthday) JSONSerialize() BirthdayJSON {
	return BirthdayJSON{birthdayJSONFields{
		ID:        b.id,
		Name:      b.name,
		Date:      b.date.String(),
		WishList:  b.wishList.String(),
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
		DeletedAt: b.deletedAt,
	}}
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

// CreateBirthdayMessage : 誕生日メッセージを生成
func (b Birthday) CreateBirthdayMessage() string {
	wishList := b.wishList.String()
	if b.wishList.IsNull() {
		wishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.name + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + wishList + " :gainings:"
}

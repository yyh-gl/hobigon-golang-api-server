package birthday

import (
	"encoding/json"
	"fmt"
)

// Birthday : 誕生日を表すドメインモデル
type Birthday struct {
	f fields
}

type fields struct {
	Name     Name     `json:"name,omitempty"`
	Date     Date     `json:"date,omitempty"`
	WishList WishList `json:"wish_list,omitempty"`
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
		fields{
			Name:     *n,
			Date:     *d,
			WishList: *wl,
		},
	}, nil
}

// Name : name のゲッター
func (b Birthday) Name() Name {
	return b.f.Name
}

// Date : date のゲッター
func (b Birthday) Date() Date {
	return b.f.Date
}

// WishList : wishList のゲッター
func (b Birthday) WishList() WishList {
	return b.f.WishList
}

// CreateBirthdayMessage : 誕生日メッセージを生成
func (b Birthday) CreateBirthdayMessage() string {
	wishList := b.f.WishList.String()
	if b.f.WishList.IsNull() {
		wishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.f.Name.String() + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + wishList + " :gainings:"
}

// MarshalJSON : Marshal用関数
// FIXME: ドメインモデル内に持ちたくないが、フィールドを公開したくもないので一旦これでいく。よりよい方法を探す
func (b Birthday) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.f)
}

// UnmarshalJSON : Unmarshal用関数
// FIXME: ドメインモデル内に持ちたくないが、フィールドを公開したくもないので一旦これでいく。よりよい方法を探す
//        テストのためにだけに用意しているので、いっそう見直したい
func (b *Birthday) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &b.f)
	if err != nil {
		return fmt.Errorf("Birthday.UnmarshalJSON()内でエラー: %w", err)
	}
	return nil
}

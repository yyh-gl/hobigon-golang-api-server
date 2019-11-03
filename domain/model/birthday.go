package model

import (
	"errors"
	"strings"
	"time"
)

type Birthday struct {
	ID        uint
	Name      string
	Date      string
	WishList  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewBirthday(
	name string,
	date time.Time,
	wishList string,
) (*Birthday, error) {
	// WishList のバリデーション
	//  -> WishList が空じゃない場合は "https://" で始まっていることをチェック
	if wishList != "" && !strings.HasPrefix(wishList, "https://") {
		return nil, errors.New("バリデーションエラー：【Birthday】WishList が \"https://\" から始まっていません")
	}

	return &Birthday{
		Name:     name,
		Date:     convertStringDate(date),
		WishList: wishList,
	}, nil
}

// time.Time 形式の日付を文字列形式（MMdd）に変換
func convertStringDate(date time.Time) string {
	return date.Format("0102")
}

// 指定 Birthday が本日のものかどうか判定
func (b Birthday) IsToday() bool {
	today := convertStringDate(time.Now())
	return b.Date == today
}

// 誕生日メッセージを生成
func (b Birthday) CreateBirthdayMessage() string {
	if b.WishList == "" {
		b.WishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.Name + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + b.WishList + " :gainings:"
}

package model

import (
	"time"
)

// TODO: ドメイン貧血症を治す
type Birthday struct {
	ID        uint
	Name      string
	Date      string
	WishList  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (b Birthday) IsToday() bool {
	today := time.Now().Format("0102")
	return b.Date == today
}

func (b Birthday) CreateBirthdayMessage() string {
	if b.WishList == "" {
		b.WishList = "Amazon の欲しい物リスト教えて！"
	}
	return "今日は *" + b.Name + "* の誕生日ンゴ > :honda:\n↓ *WishList* ↓\n:gainings: " + b.WishList + " :gainings:"
}

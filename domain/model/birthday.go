package model

import (
	"time"
)

// TODO: ドメイン貧血症を治す
// TODO: JSON タグをドメインモデルではなく、ハンドラー層に定義した構造体に定義するように修正する
// TODO: gorm タグを削除
type Birthday struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"name;not null"`
	Date      string `gorm:"date;not null"`
	WishList  string `gorm:"wish_list"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (b Birthday) TableName() string {
	return "birthdays"
}

func (b Birthday) IsToday() bool {
	today := time.Now().Format("0102")
	return b.Date == today
}

func (b Birthday) CreateBirthdayMessage() string {
	return "今日は *" + b.Name + "* の誕生日だっぴ > :honda:\n:gainings: " + b.WishList + " :gainings:"
}

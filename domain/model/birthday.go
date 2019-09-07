package model

import (
	"time"
)

// TODO: ドメイン貧血症を治す
type Birthday struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	Name      string     `json:"name,omitempty" gorm:"name;not null"`
	Date      string     `json:"date,omitempty" gorm:"date;not null"`
	WishList  string     `json:"wish_list,omitempty" gorm:"wish_list"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
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

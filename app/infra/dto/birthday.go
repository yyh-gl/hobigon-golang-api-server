package dto

import (
	"time"
)

// BirthdayDTO : 誕生日用のDTO
type BirthdayDTO struct {
	Name      string `gorm:"not null"`
	Date      string `gorm:"index;not null"`
	WishList  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// TableName : DB アクセスにおける対応テーブル名
func (b BirthdayDTO) TableName() string {
	return "birthdays"
}

package model

import "time"

type BirthdayDTO struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"name;not null"`
	Date      string `gorm:"date;not null"`
	WishList  string `gorm:"wish_list"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (b BirthdayDTO) TableName() string {
	return "birthdays"
}

package model

import (
	"time"
)

type Blog struct {
	// TODO: gorm.Model を使用する
	ID        uint `json:"id,omitempty",gorm:"primary_key"`
	Title     string `json:"title,omitempty",gorm:"title"`
	Count     *int64  `json:"count,omitempty"gorm:"count"`
	CreatedAt *time.Time `json:"created_at,omitempty"gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"gorm:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"gorm:"deleted_at",sql:"index"`
}

func (b Blog) TableName() string {
	return "blog_posts"
}

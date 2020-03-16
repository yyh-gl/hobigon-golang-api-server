package dto

import (
	"time"
)

// BlogDTO : ブログ用のDTO
type BlogDTO struct {
	Title     string `gorm:"primary_key;not null"`
	Count     int    `gorm:"default:0;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// TableName : DBアクセスにおける対応テーブル名
func (b BlogDTO) TableName() string {
	return "blog_posts"
}

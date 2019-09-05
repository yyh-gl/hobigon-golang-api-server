package model

import "time"

type Blog struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	Title     string     `json:"title,omitempty" gorm:"title;unique;not null"`
	Count     *int       `json:"count,omitempty" gorm:"count;default:0;not null"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

func (b Blog) TableName() string {
	return "blog_posts"
}

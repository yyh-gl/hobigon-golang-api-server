package entity

import "time"

// TODO: ドメイン貧血症を治す
type Blog struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Title     string `gorm:"title;unique;not null"`
	Count     *int   `gorm:"count;default:0;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (b Blog) TableName() string {
	return "blog_posts"
}

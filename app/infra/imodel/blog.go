package imodel

import (
	"context"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
)

// BlogDTO : ブログ用の DTO
type BlogDTO struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Title     string `gorm:"unique;not null"`
	Count     *int   `gorm:"default:0;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// TableName : DB アクセスにおける対応テーブル名
func (b BlogDTO) TableName() string {
	return "blog_posts"
}

// ConvertToDomainModel : ドメインモデルに変換
func (b BlogDTO) ConvertToDomainModel(ctx context.Context) *blog.Blog {
	return blog.NewBlogWithFullParams(
		b.ID, b.Title, b.Count, b.CreatedAt, b.UpdatedAt, b.DeletedAt,
	)
}

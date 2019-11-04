package imodel

import (
	"context"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
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
func (b BlogDTO) ConvertToDomainModel(ctx context.Context) *entity.Blog {
	return entity.NewBlogWithFullParams(
		b.ID, b.Title, b.Count, b.CreatedAt, b.UpdatedAt, b.DeletedAt,
	)
}

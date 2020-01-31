package imodel

import (
	"context"
	"fmt"
	"time"

	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
)

// BlogDTO : ブログ用の DTO
type BlogDTO struct {
	Title     string `gorm:"primary_key;not null"`
	Count     int    `gorm:"default:0;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// TableName : DB アクセスにおける対応テーブル名
func (b BlogDTO) TableName() string {
	return "blog_posts"
}

// ConvertToDomainModel : ドメインモデルに変換
func (b BlogDTO) ConvertToDomainModel(ctx context.Context) (*model.Blog, error) {
	blog, err := model.NewBlogWithFullParams(b.Title, b.Count)
	if err != nil {
		return nil, fmt.Errorf("model.NewBlogWithFullParams()でエラー: %w", err)
	}
	return blog, nil
}

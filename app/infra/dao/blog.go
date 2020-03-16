package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

type blog struct {
	db *db.DB
}

// NewBlog : Blog用リポジトリを取得
func NewBlog(db *db.DB) repository.Blog {
	return &blog{
		db: db,
	}
}

// Create : ブログ情報を新規作成
func (b blog) Create(ctx context.Context, blog model.Blog) (*model.Blog, error) {
	// BlogモデルをDTOに変換
	blogDTO := dto.BlogDTO{
		Title: blog.Title().String(),
		Count: blog.Count().Int(),
	}

	err := b.db.Create(&blogDTO).Error
	if err != nil {
		return nil, fmt.Errorf("gorm.Create(blog)内でのエラー: %w", err)
	}

	createdBlog := model.ConvertToDomainModel(ctx, blogDTO)
	return createdBlog, nil
}

// FindByTitle : タイトルからブログ情報を1件取得
func (b blog) FindByTitle(ctx context.Context, title string) (*model.Blog, error) {
	blogDTO := dto.BlogDTO{}
	err := b.db.First(&blogDTO, "title=?", title).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("gorm.First(blog)内でのエラー: %w", err)
	}

	blog := model.ConvertToDomainModel(ctx, blogDTO)
	return blog, nil
}

// Update : ブログ情報を1件更新
func (b blog) Update(ctx context.Context, blog model.Blog) (*model.Blog, error) {
	// BlogモデルをDTOに変換
	blogDTO := dto.BlogDTO{
		Title: blog.Title().String(),
		Count: blog.Count().Int(),
	}

	err := b.db.Save(&blogDTO).Error
	if err != nil {
		return nil, fmt.Errorf("gorm.Save(blog)内でのエラー: %w", err)
	}

	updatedBlog := model.ConvertToDomainModel(ctx, blogDTO)
	return updatedBlog, nil
}

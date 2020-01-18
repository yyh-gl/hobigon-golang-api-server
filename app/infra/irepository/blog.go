package irepository

import (
	"context"

	"github.com/pkg/errors"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/imodel"
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

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : ブログ情報を新規作成
func (b blog) Create(ctx context.Context, blog model.Blog) (*model.Blog, error) {
	// Blog モデル を DTO に変換
	blogDTO := imodel.BlogDTO{
		Title: blog.Title(),
		Count: blog.Count(),
	}

	err := b.db.Create(&blogDTO).Error
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Create(blog)内でのエラー")
	}

	createdBlog := blogDTO.ConvertToDomainModel(ctx)
	return createdBlog, nil
}

//////////////////////////////////////////////////
// SelectByTitle
//////////////////////////////////////////////////

// SelectByTitle : タイトルからブログ情報を1件取得
func (b blog) SelectByTitle(ctx context.Context, title string) (*model.Blog, error) {
	blogDTO := imodel.BlogDTO{}
	err := b.db.First(&blogDTO, "title=?", title).Error
	if err != nil {
		return nil, errors.Wrap(err, "gorm.First(blog)内でのエラー")
	}

	blog := blogDTO.ConvertToDomainModel(ctx)
	return blog, nil
}

//////////////////////////////////////////////////
// Update
//////////////////////////////////////////////////

// Update : ブログ情報を1件更新
func (b blog) Update(ctx context.Context, blog model.Blog) (*model.Blog, error) {
	// Blog モデル を DTO に変換
	blogDTO := imodel.BlogDTO{
		ID:    blog.ID(),
		Title: blog.Title(),
		Count: blog.Count(),
	}

	err := b.db.Save(&blogDTO).Error
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Save(blog)内でのエラー")
	}

	updatedBlog := blogDTO.ConvertToDomainModel(ctx)
	return updatedBlog, nil
}

package irepository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/db"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/blog"

	"github.com/pkg/errors"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/imodel"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

//////////////////////////////////////////////////
// NewBlogRepository
//////////////////////////////////////////////////

type blogRepository struct {
	db *db.DB
}

// NewBlogRepository : ブログ用のリポジトリを取得
func NewBlogRepository(db *db.DB) repository.BlogRepository {
	return &blogRepository{
		db: db,
	}
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : ブログ情報を新規作成
func (br blogRepository) Create(ctx context.Context, blog blog.Blog) (*blog.Blog, error) {
	// Blog モデル を DTO に変換
	blogDTO := imodel.BlogDTO{
		Title: blog.Title(),
		Count: blog.Count(),
	}

	err := br.db.Create(&blogDTO).Error
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
func (br blogRepository) SelectByTitle(ctx context.Context, title string) (*blog.Blog, error) {
	blogDTO := imodel.BlogDTO{}
	err := br.db.First(&blogDTO, "title=?", title).Error
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
func (br blogRepository) Update(ctx context.Context, blog blog.Blog) (*blog.Blog, error) {
	// Blog モデル を DTO に変換
	blogDTO := imodel.BlogDTO{
		ID:    blog.ID(),
		Title: blog.Title(),
		Count: blog.Count(),
	}

	err := br.db.Save(&blogDTO).Error
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Save(blog)内でのエラー")
	}

	updatedBlog := blogDTO.ConvertToDomainModel(ctx)
	return updatedBlog, nil
}

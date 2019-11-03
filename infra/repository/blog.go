package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

//////////////////////////////////////////////////
// NewBlogRepository
//////////////////////////////////////////////////

type blogRepository struct {
	db *gorm.DB
}

// NewBlogRepository : ブログ用のリポジトリを取得
func NewBlogRepository() repository.BlogRepository {
	return &blogRepository{
		db: app.DB,
	}
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : ブログ情報を新規作成
func (br blogRepository) Create(ctx context.Context, blog entity.Blog) (entity.Blog, error) {
	err := br.db.Create(&blog).Error
	if err != nil {
		return entity.Blog{}, err
	}
	return blog, nil
}

//////////////////////////////////////////////////
// SelectByTitle
//////////////////////////////////////////////////

// SelectByTitle : タイトルからブログ情報を1件取得
func (br blogRepository) SelectByTitle(ctx context.Context, title string) (blog entity.Blog, err error) {
	err = br.db.First(&blog, "title=?", title).Error
	if err != nil {
		return entity.Blog{}, err
	}
	return blog, nil
}

//////////////////////////////////////////////////
// Update
//////////////////////////////////////////////////

// Update : ブログ情報を1件更新
func (br blogRepository) Update(ctx context.Context, blog entity.Blog) (entity.Blog, error) {
	err := br.db.Save(&blog).Error
	if err != nil {
		return entity.Blog{}, err
	}
	return blog, nil
}

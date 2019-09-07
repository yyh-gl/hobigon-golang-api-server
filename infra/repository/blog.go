package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type blogRepository struct{}

func NewBlogRepository() repository.BlogRepository {
	return &blogRepository{}
}

func (br blogRepository) Create(ctx context.Context, blog model.Blog) (model.Blog, error) {
	db := ctx.Value("db").(*gorm.DB)
	err := db.Create(&blog).Error

	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

func (br blogRepository) SelectByTitle(ctx context.Context, title string) (blog model.Blog, err error) {
	db := ctx.Value("db").(*gorm.DB)
	err = db.First(&blog, "title=?", title).Error
	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

func (br blogRepository) Update(ctx context.Context, blog model.Blog) (model.Blog, error) {
	db := ctx.Value("db").(*gorm.DB)
	err := db.Save(&blog).Error
	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository() repository.BlogRepository {
	return &blogRepository{
		db: app.DB,
	}
}

func (br blogRepository) Create(blog model.Blog) (model.Blog, error) {
	err := br.db.Create(&blog).Error
	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

func (br blogRepository) SelectByTitle(title string) (blog model.Blog, err error) {
	err = br.db.First(&blog, "title=?", title).Error
	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

func (br blogRepository) Update(blog model.Blog) (model.Blog, error) {
	err := br.db.Save(&blog).Error
	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

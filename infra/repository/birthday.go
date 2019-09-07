package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type birthdayRepository struct {
	db *gorm.DB
}

func NewBirthdayRepository() repository.BirthdayRepository {
	return &birthdayRepository{
		db: app.DB,
	}
}

func (br birthdayRepository) SelectByDate(date string) (birthday model.Birthday, err error) {
	err = br.db.First(&birthday, "date=?", date).Error
	if err != nil {
		return model.Birthday{}, err
	}
	return birthday, nil
}

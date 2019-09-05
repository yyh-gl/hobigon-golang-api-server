package repository

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type birthdayRepository struct {}

func NewBirthdayRepository() repository.BirthdayRepository {
	return &birthdayRepository{}
}

// TODO: ドメイン層への依存をなくす
func (br birthdayRepository) SelectByDate(ctx context.Context, date string) (birthday model.Birthday, err error) {
	db := ctx.Value("db").(*gorm.DB)
	err = db.First(&birthday, "date=?", date).Error
	if err != nil {
		return model.Birthday{}, err
	}
	return birthday, nil
}

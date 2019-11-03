package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

//////////////////////////////////////////////////
// NewBirthdayRepository
//////////////////////////////////////////////////

type birthdayRepository struct {
	db *gorm.DB
}

// NewBirthdayRepository : 誕生日用のリポジトリを取得
func NewBirthdayRepository() repository.BirthdayRepository {
	return &birthdayRepository{
		db: app.DB,
	}
}

//////////////////////////////////////////////////
// SelectByDate
//////////////////////////////////////////////////

// SelectByDate : 日付から誕生日を1件取得
func (br birthdayRepository) SelectByDate(ctx context.Context, date string) (birthday model.Birthday, err error) {
	err = br.db.First(&birthday, "date=?", date).Error
	if err != nil {
		return model.Birthday{}, err
	}
	return birthday, nil
}

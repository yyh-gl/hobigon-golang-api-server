package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
	infraModel "github.com/yyh-gl/hobigon-golang-api-server/infra/model"
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
	// Birthday の DTO を用意
	birthdayDTO := infraModel.BirthdayDTO{}

	// date 指定で誕生日情報を取得
	err = br.db.First(&birthdayDTO, "date=?", date).Error
	if err != nil {
		return model.Birthday{}, err
	}

	// DTO を ドメインモデルに変換
	birthday = model.Birthday(birthdayDTO)

	return birthday, nil
}

package repository

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
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
// Create
//////////////////////////////////////////////////

// Create : 誕生日データを新規作成
func (br birthdayRepository) Create(ctx context.Context, birthday entity.Birthday) (*entity.Birthday, error) {
	// Birthday モデル を DTO に変換
	birthdayDTO := infraModel.BirthdayDTO{
		Name:     birthday.Name(),
		Date:     birthday.Date().String(),
		WishList: birthday.WishList().String(),
	}

	// date 指定で誕生日情報を取得
	err := br.db.Create(&birthdayDTO).Error
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Create()内でのエラー")
	}

	// DTO を ドメインモデルに変換
	createdBirthday, err := birthdayDTO.ConvertToDomainModel()
	if err != nil {
		return nil, errors.Wrap(err, "birthdayDTO.ConvertToDomainModel()内でのエラー")
	}
	return createdBirthday, nil
}

//////////////////////////////////////////////////
// SelectByDate
//////////////////////////////////////////////////

// SelectByDate : 日付から誕生日を1件取得
func (br birthdayRepository) SelectByDate(ctx context.Context, date string) (*entity.Birthday, error) {
	// Birthday の DTO を用意
	birthdayDTO := infraModel.BirthdayDTO{}

	// date 指定で誕生日情報を取得
	err := br.db.First(&birthdayDTO, "date=?", date).Error
	if err != nil {
		return nil, err
	}

	// DTO を ドメインモデルに変換
	birthday, err := birthdayDTO.ConvertToDomainModel()
	if err != nil {
		return nil, errors.Wrap(err, "birthdayDTO.ConvertToDomainModel()内でのエラー")
	}
	return birthday, nil
}

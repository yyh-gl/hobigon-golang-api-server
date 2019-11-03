package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"

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
// Create
//////////////////////////////////////////////////

// Create : 誕生日データを新規作成
func (br birthdayRepository) Create(ctx context.Context, birthday model.Birthday) (*model.Birthday, error) {
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
	month, err := strconv.Atoi(birthdayDTO.Date[0:2])
	if err != nil {
		return nil, errors.Wrap(err, "hour取得におけるstrconv.Atoi()内でのエラー")
	}
	day, err := strconv.Atoi(birthdayDTO.Date[2:4])
	if err != nil {
		return nil, errors.Wrap(err, "hour取得におけるstrconv.Atoi()内でのエラー")
	}
	d := time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.Local)
	createdBirthday, err := model.NewBirthday(birthdayDTO.Name, d, birthdayDTO.WishList)
	if err != nil {
		return nil, errors.Wrap(err, "NewBirthday()内でのエラー")
	}
	createdBirthday.SetID(birthdayDTO.ID)
	createdBirthday.SetCreatedAt(birthdayDTO.CreatedAt)
	createdBirthday.SetUpdatedAt(birthdayDTO.UpdatedAt)
	createdBirthday.SetDeletedAt(birthdayDTO.DeletedAt)

	return createdBirthday, nil
}

//////////////////////////////////////////////////
// SelectByDate
//////////////////////////////////////////////////

// SelectByDate : 日付から誕生日を1件取得
func (br birthdayRepository) SelectByDate(ctx context.Context, date string) (*model.Birthday, error) {
	// Birthday の DTO を用意
	birthdayDTO := infraModel.BirthdayDTO{}

	// date 指定で誕生日情報を取得
	err := br.db.First(&birthdayDTO, "date=?", date).Error
	if err != nil {
		return nil, err
	}

	// DTO を ドメインモデルに変換
	month, err := strconv.Atoi(birthdayDTO.Date[0:2])
	if err != nil {
		return nil, errors.Wrap(err, "hour取得におけるstrconv.Atoi()内でのエラー")
	}
	day, err := strconv.Atoi(birthdayDTO.Date[2:4])
	if err != nil {
		return nil, errors.Wrap(err, "hour取得におけるstrconv.Atoi()内でのエラー")
	}
	d := time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.Local)
	birthday, err := model.NewBirthday(birthdayDTO.Name, d, birthdayDTO.WishList)
	if err != nil {
		return nil, errors.Wrap(err, "NewBirthday()内でのエラー")
	}
	birthday.SetID(birthdayDTO.ID)
	birthday.SetCreatedAt(birthdayDTO.CreatedAt)
	birthday.SetUpdatedAt(birthdayDTO.UpdatedAt)
	birthday.SetDeletedAt(birthdayDTO.DeletedAt)

	return birthday, nil
}

package irepository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/imodel"
)

//////////////////////////////////////////////////
// NewBirthdayRepository
//////////////////////////////////////////////////

type birthdayRepository struct {
	db *db.DB
}

// NewBirthdayRepository : 誕生日用のリポジトリを取得
func NewBirthdayRepository(db *db.DB) repository.BirthdayRepository {
	return &birthdayRepository{
		db: db,
	}
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : 誕生日データを新規作成
func (br birthdayRepository) Create(ctx context.Context, birthday birthday.Birthday) (*birthday.Birthday, error) {
	// Birthday モデル を DTO に変換
	birthdayDTO := imodel.BirthdayDTO{
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
	createdBirthday, err := birthdayDTO.ConvertToDomainModel(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "birthdayDTO.ConvertToDomainModel()内でのエラー")
	}
	return createdBirthday, nil
}

//////////////////////////////////////////////////
// SelectByDate
//////////////////////////////////////////////////

// SelectByDate : 日付から誕生日を1件取得
func (br birthdayRepository) SelectByDate(ctx context.Context, date string) (*birthday.Birthday, error) {
	// Birthday の DTO を用意
	birthdayDTO := imodel.BirthdayDTO{}

	// date 指定で誕生日情報を取得
	err := br.db.First(&birthdayDTO, "date=?", date).Error
	if err != nil {
		return nil, err
	}

	// DTO を ドメインモデルに変換
	birthday, err := birthdayDTO.ConvertToDomainModel(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "birthdayDTO.ConvertToDomainModel()内でのエラー")
	}
	return birthday, nil
}

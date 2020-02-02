package dao

import (
	"context"
	"fmt"

	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

type birthday struct {
	db *db.DB
}

// NewBirthday : 誕生日用のリポジトリを取得
func NewBirthday(db *db.DB) repository.Birthday {
	return &birthday{
		db: db,
	}
}

// Create : 誕生日データを新規作成
func (b birthday) Create(ctx context.Context, birthday model.Birthday) (*model.Birthday, error) {
	// Birthday モデル を DTO に変換
	birthdayDTO := dto.BirthdayDTO{
		Name:     birthday.Name().String(),
		Date:     birthday.Date().String(),
		WishList: birthday.WishList().String(),
	}

	// date 指定で誕生日情報を取得
	err := b.db.Create(&birthdayDTO).Error
	if err != nil {
		return nil, fmt.Errorf("gorm.Create()内でのエラー: %w", err)
	}

	// DTO を ドメインモデルに変換
	createdBirthday, err := birthdayDTO.ConvertToDomainModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("birthdayDTO.ConvertToDomainModel()内でのエラー: %w", err)
	}
	return createdBirthday, nil
}

// SelectByDate : 日付から誕生日を1件取得
func (b birthday) SelectByDate(ctx context.Context, date string) (*model.Birthday, error) {
	// Birthday の DTO を用意
	birthdayDTO := dto.BirthdayDTO{}

	// date 指定で誕生日情報を取得
	err := b.db.First(&birthdayDTO, "date=?", date).Error
	if err != nil {
		return nil, fmt.Errorf("gorm.First()内でのエラー: %w", err)
	}

	// DTO を ドメインモデルに変換
	birthday, err := birthdayDTO.ConvertToDomainModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("birthdayDTO.ConvertToDomainModel()内でのエラー: %w", err)
	}
	return birthday, nil
}

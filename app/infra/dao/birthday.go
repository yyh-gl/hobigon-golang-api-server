package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
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

// FindAllByDate : 指定日付の誕生日データを全件取得
func (b birthday) FindAllByDate(ctx context.Context, date string) (*model.BirthdayList, error) {
	// date 指定で誕生日情報を取得
	var birthdayListDTO dto.BirthdayListDTO
	err := b.db.Where("date=?", date).Find(&birthdayListDTO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRecordNotFound
		}
		return nil, fmt.Errorf("gorm.Where().Find()内でのエラー: %w", err)
	}

	// DTO を ドメインモデルに変換
	birthdayList, err := birthdayListDTO.ConvertToDomainModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("birthdayListDTO.ConvertToDomainModel()内でのエラー: %w", err)
	}
	return birthdayList, nil
}

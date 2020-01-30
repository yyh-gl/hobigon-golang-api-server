package imodel

import (
	"context"
	"fmt"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

// BirthdayDTO : 誕生日用の DTO
type BirthdayDTO struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"not null"`
	Date      string `gorm:"not null"`
	WishList  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// TableName : DB アクセスにおける対応テーブル名
func (b BirthdayDTO) TableName() string {
	return "birthdays"
}

// ConvertToDomainModel : ドメインモデルに変換
func (b BirthdayDTO) ConvertToDomainModel(ctx context.Context) (*birthday.Birthday, error) {
	// Birthday モデルを取得
	domainModelBirthday, err := birthday.NewBirthdayWithFullParams(
		b.ID, b.Name, b.Date, b.WishList, b.CreatedAt, b.UpdatedAt, b.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("birthday.NewBirthdayWithFullParams()でエラー: %w", err)
	}

	return domainModelBirthday, nil
}

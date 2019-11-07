package imodel

import (
	"context"
	"strconv"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/birthday"

	"github.com/pkg/errors"
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
	// time.Time 型の日付情報を取得
	month, err := strconv.Atoi(b.Date[0:2])
	if err != nil {
		return nil, errors.Wrap(err, "hour取得におけるstrconv.Atoi()内でのエラー")
	}
	day, err := strconv.Atoi(b.Date[2:4])
	if err != nil {
		return nil, errors.Wrap(err, "hour取得におけるstrconv.Atoi()内でのエラー")
	}
	date := time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.Local)

	// Birthday モデルを取得
	birthday, err := birthday.NewBirthdayWithFullParams(
		b.ID, b.Name, date, b.WishList, b.CreatedAt, b.UpdatedAt, b.DeletedAt,
	)
	return birthday, nil
}

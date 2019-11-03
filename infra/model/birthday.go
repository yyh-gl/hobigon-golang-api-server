package model

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

type BirthdayDTO struct {
	ID        uint   `gorm:"primary_key;AUTO_INCREMENT"`
	Name      string `gorm:"name;not null"`
	Date      string `gorm:"date;not null"`
	WishList  string `gorm:"wish_list"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (b BirthdayDTO) TableName() string {
	return "birthdays"
}

func (b BirthdayDTO) ConvertToDomainModel() (*entity.Birthday, error) {
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
	birthday, err := entity.NewBirthdayWithFullParams(
		b.ID, b.Name, date, b.WishList, b.CreatedAt, b.UpdatedAt, b.DeletedAt,
	)
	return birthday, nil
}

package usecase

import (
	"context"
	"fmt"
	"time"

	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
)

// Birthday : Birthday用ユースケースのインターフェース
type Birthday interface {
	Create(ctx context.Context, name string, date time.Time, wishList string) (*model.Birthday, error)
}

type birthday struct {
	r repository.Birthday
}

// NewBirthday : Birthday用ユースケースを取得
func NewBirthday(
	r repository.Birthday,
) Birthday {
	return &birthday{
		r: r,
	}
}

// Create : 誕生日データを新規作成
func (b birthday) Create(ctx context.Context, name string, date time.Time, wishList string) (*model.Birthday, error) {
	// 新しい Birthday データを作成
	newBirthday, err := model.NewBirthday(name, date, wishList)
	if err != nil {
		return nil, fmt.Errorf("birthday.New()内でエラー: %w", err)
	}

	createdBirthday, err := b.r.Create(ctx, *newBirthday)
	if err != nil {
		return nil, fmt.Errorf("birthdayRepository.Create()内でエラー: %w", err)
	}
	return createdBirthday, nil
}

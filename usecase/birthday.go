package usecase

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

//////////////////////////////////////////////////
// NewBirthdayUseCase
//////////////////////////////////////////////////

// BirthdayUseCase : 誕生日用のユースケースインターフェース
type BirthdayUseCase interface {
	Create(ctx context.Context, name string, date time.Time, wishList string) (*entity.Birthday, error)
}

type birthdayUseCase struct {
	br repository.BirthdayRepository
}

// NewBirthdayUseCase : 通知用のユースケースを取得
func NewBirthdayUseCase(
	br repository.BirthdayRepository,
) BirthdayUseCase {
	return &birthdayUseCase{
		br: br,
	}
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : 誕生日データを新規作成
func (bu birthdayUseCase) Create(ctx context.Context, name string, date time.Time, wishList string) (*entity.Birthday, error) {
	// 新しい Birthday データを作成
	birthday, err := entity.NewBirthday(name, date, wishList)
	if err != nil {
		return nil, errors.Wrap(err, "NewBirthday()内でのエラー")
	}

	createdBirthday, err := bu.br.Create(ctx, *birthday)
	if err != nil {
		return nil, errors.Wrap(err, "birthdayRepository.Create()内でのエラー")
	}
	return createdBirthday, nil
}

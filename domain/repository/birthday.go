package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

// BirthdayRepository : 誕生日用のリポジトリインターフェース
type BirthdayRepository interface {
	Create(ctx context.Context, birthday entity.Birthday) (*entity.Birthday, error)
	SelectByDate(ctx context.Context, date string) (*entity.Birthday, error)
}

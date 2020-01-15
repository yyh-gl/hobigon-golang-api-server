package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

// BirthdayRepository : 誕生日用のリポジトリインターフェース
type BirthdayRepository interface {
	Create(ctx context.Context, birthday birthday.Birthday) (*birthday.Birthday, error)
	SelectByDate(ctx context.Context, date string) (*birthday.Birthday, error)
}

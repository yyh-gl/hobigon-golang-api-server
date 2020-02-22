package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

// Birthday : Birthday用リポジトリのインターフェース
type Birthday interface {
	Create(ctx context.Context, birthday birthday.Birthday) (*birthday.Birthday, error)
	FindAllByDate(ctx context.Context, date string) (birthday.BirthdayList, error)
}

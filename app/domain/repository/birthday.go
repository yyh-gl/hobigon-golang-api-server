package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

// Birthday : Birthday用リポジトリのインターフェース
type Birthday interface {
	Create(ctx context.Context, birthday birthday.Birthday) (*birthday.Birthday, error)
	SelectByDate(ctx context.Context, date string) (*birthday.Birthday, error)
}

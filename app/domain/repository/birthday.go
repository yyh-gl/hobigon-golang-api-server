package repository

import (
	"context"
	"errors"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

// ErrBirthdayRecordNotFound : DBに該当Birthday情報が存在しないエラー
var ErrBirthdayRecordNotFound = errors.New("record of birthday is not found")

// Birthday : Birthday用リポジトリのインターフェース
type Birthday interface {
	Create(ctx context.Context, birthday birthday.Birthday) (birthday.Birthday, error)
	FindAllByDate(ctx context.Context, date string) (birthday.BirthdayList, error)
}

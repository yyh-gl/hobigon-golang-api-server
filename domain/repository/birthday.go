package repository

import (
	"context"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type BirthdayRepository interface {
	SelectByDate(ctx context.Context, date string) (model.Birthday, error)
}

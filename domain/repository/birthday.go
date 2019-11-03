package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type BirthdayRepository interface {
	Create(ctx context.Context, birthday model.Birthday) (*model.Birthday, error)
	SelectByDate(ctx context.Context, date string) (*model.Birthday, error)
}

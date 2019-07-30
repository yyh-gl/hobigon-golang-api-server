package repository

import (
	"context"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type BlogRepository interface {
	SelectByTitle(ctx context.Context, title string) (model.Blog, error)
}

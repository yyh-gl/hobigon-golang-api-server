package repository

import (
	"context"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type BlogRepository interface {
	Create(ctx context.Context, blog model.Blog) (model.Blog, error)
	SelectByTitle(ctx context.Context, title string) (model.Blog, error)
	Update(ctx context.Context, blog model.Blog) (model.Blog, error)
}

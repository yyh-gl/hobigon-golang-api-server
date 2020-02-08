package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
)

// Blog : Blog用リポジトリのインターフェース
type Blog interface {
	Create(ctx context.Context, blog blog.Blog) (*blog.Blog, error)
	FindByTitle(ctx context.Context, title string) (*blog.Blog, error)
	Update(ctx context.Context, blog blog.Blog) (*blog.Blog, error)
}

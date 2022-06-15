package repository

import (
	"context"
	"errors"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
)

// ErrBlogRecordNotFound : DBに該当Blog情報が存在しないエラー
var ErrBlogRecordNotFound = errors.New("record of blog is not found")

// Blog : Blog用リポジトリのインターフェース
type Blog interface {
	Create(ctx context.Context, blog blog.Blog) (blog.Blog, error)
	FindByTitle(ctx context.Context, title string) (blog.Blog, error)
	Update(ctx context.Context, blog blog.Blog) (blog.Blog, error)
}

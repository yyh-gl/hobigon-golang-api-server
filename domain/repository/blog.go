package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

// BlogRepository : ブログ用のリポジトリインターフェース
type BlogRepository interface {
	Create(ctx context.Context, blog entity.Blog) (entity.Blog, error)
	SelectByTitle(ctx context.Context, title string) (entity.Blog, error)
	Update(ctx context.Context, blog entity.Blog) (entity.Blog, error)
}

package repository

import (
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// BlogRepository : ブログ用のリポジトリインターフェース
type BlogRepository interface {
	Create(blog model.Blog) (model.Blog, error)
	SelectByTitle(title string) (model.Blog, error)
	Update(blog model.Blog) (model.Blog, error)
}

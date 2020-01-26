package usecase

import (
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
)

// ErrBlogNotFound : 指定ブログが存在しないエラー
var ErrBlogNotFound = repository.ErrRecordNotFound

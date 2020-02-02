package usecase

import (
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
)

// ErrBlogNotFound : 該当Blogが存在しないエラー
var ErrBlogNotFound = repository.ErrRecordNotFound

// ErrBirthdayNotFound : 該当Birthdayが存在しないエラー
var ErrBirthdayNotFound = repository.ErrRecordNotFound

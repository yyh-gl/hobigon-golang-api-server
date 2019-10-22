package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

type CreateBlogParams struct {
	Title string
}

// CreateBlogUseCase はブログデータを新規で作成
func CreateBlogUseCase(ctx context.Context, params CreateBlogParams) (*model.Blog, error) {
	blogRepository := repository.NewBlogRepository()

	blog := model.Blog{
		Title: params.Title,
	}
	blog, err := blogRepository.Create(blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Create()内でのエラー")
	}

	return &blog, nil
}

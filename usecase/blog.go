package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

//////////////////////////////////////////////////
// CreateBlog
//////////////////////////////////////////////////

// CreateBlogParams はブログデータ作成に必要なパラメータを受け取るための構造体
type CreateBlogParams struct {
	Title string
}

// CreateBlogUseCase はブログデータを新規で作成
func CreateBlogUseCase(ctx context.Context, title string) (*model.Blog, error) {
	blogRepository := repository.NewBlogRepository()

	blog := model.Blog{
		Title: title,
	}
	blog, err := blogRepository.Create(blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Create()内でのエラー")
	}

	return &blog, nil
}

//////////////////////////////////////////////////
// GetBlog
//////////////////////////////////////////////////

// GetBlogUseCase はブログデータを1件取得
func GetBlogUseCase(ctx context.Context, title string) (*model.Blog, error) {
	blogRepository := repository.NewBlogRepository()

	blog, err := blogRepository.SelectByTitle(title)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, err
		default:
			return nil, errors.Wrap(err, "blogRepository.SelectByTitle()内でのエラー")
		}
	}

	return &blog, nil
}

//////////////////////////////////////////////////
// LikeBlog
//////////////////////////////////////////////////

// LikeBlogUseCase は指定ブログにいいねをプラス1
func LikeBlogUseCase(ctx context.Context, title string) (*model.Blog, error) {
	blogRepository := repository.NewBlogRepository()
	slackGateway := gateway.NewSlackGateway()

	blog, err := blogRepository.SelectByTitle(title)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, err
		default:
			return nil, errors.Wrap(err, "blogRepository.SelectByTitle()内でのエラー")
		}
	}

	// Count をプラス1
	addedCount := *blog.Count + 1
	blog.Count = &addedCount
	blog, err = blogRepository.Update(blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Update()内でのエラー")
	}

	// Slack に通知
	err = slackGateway.SendLikeNotify(blog)
	if err != nil {
		return nil, errors.Wrap(err, "slackGateway.SendLikeNotify()内でのエラー")
	}

	return &blog, nil
}

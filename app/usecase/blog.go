package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
)

// Blog : Blog用ユースケースのインターフェース
type Blog interface {
	Create(context.Context, string) (*model.Blog, error)
	Show(context.Context, string) (*model.Blog, error)
	Like(context.Context, string) (*model.Blog, error)
}

type blog struct {
	r  repository.Blog
	sg gateway.SlackGateway
}

// NewBlog : Blog用ユースケースを取得
func NewBlog(r repository.Blog, sg gateway.SlackGateway) Blog {
	return &blog{
		r:  r,
		sg: sg,
	}
}

// Create : ブログ情報を新規作成
func (b blog) Create(ctx context.Context, title string) (*model.Blog, error) {
	blog := model.NewBlog(title)
	createdBlog, err := b.r.Create(ctx, *blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Create()内でのエラー")
	}

	return createdBlog, nil
}

// Show : ブログ情報を1件取得
func (b blog) Show(ctx context.Context, title string) (*model.Blog, error) {
	blog, err := b.r.SelectByTitle(ctx, title)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, err
		default:
			return nil, errors.Wrap(err, "blogRepository.SelectByTitle()内でのエラー")
		}
	}

	return blog, nil
}

// Like : 指定ブログにいいねをプラス1
func (b blog) Like(ctx context.Context, title string) (*model.Blog, error) {
	blog, err := b.r.SelectByTitle(ctx, title)
	if err != nil {
		switch err.Error() {
		case "record not found":
			return nil, err
		default:
			return nil, errors.Wrap(err, "blogRepository.SelectByTitle()内でのエラー")
		}
	}

	// Count をプラス1
	blog.CountUp()
	blog, err = b.r.Update(ctx, *blog)
	if err != nil {
		return nil, errors.Wrap(err, "blogRepository.Update()内でのエラー")
	}

	// Slack に通知
	err = b.sg.SendLikeNotify(ctx, *blog)
	if err != nil {
		return nil, errors.Wrap(err, "slackGateway.SendLikeNotify()内でのエラー")
	}

	return blog, nil
}

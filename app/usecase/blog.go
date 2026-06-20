package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

// ErrBlogNotFound : 該当Blogが存在しないエラー
var ErrBlogNotFound = errors.New("blog is not found")

// Blog : Blog用ユースケースのインターフェース
type Blog interface {
	Create(context.Context, string) (model.Blog, error)
	Show(context.Context, string) (model.Blog, error)
	Like(context.Context, string, bool) (model.Blog, error)
}

type blog struct {
	r  repository.Blog
	sg gateway.Slack
}

// NewBlog : Blog用ユースケースを取得
func NewBlog(r repository.Blog, sg gateway.Slack) Blog {
	return &blog{
		r:  r,
		sg: sg,
	}
}

// Create : ブログ情報を新規作成
func (b blog) Create(ctx context.Context, title string) (model.Blog, error) {
	blog, err := model.NewBlog(title)
	if err != nil {
		return model.Blog{}, fmt.Errorf("model.NewBlog()内でエラー: %w", err)
	}

	createdBlog, err := b.r.Create(ctx, blog)
	if err != nil {
		return model.Blog{}, fmt.Errorf("blogRepository.Create()内でのエラー: %w", err)
	}
	return createdBlog, nil
}

// Show : ブログ情報を1件取得
func (b blog) Show(ctx context.Context, title string) (model.Blog, error) {
	blog, err := b.r.FindByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, repository.ErrBlogRecordNotFound) {
			return model.Blog{}, ErrBlogNotFound
		}
		return model.Blog{}, fmt.Errorf("blogRepository.FindByTitle()内でのエラー: %w", err)
	}
	return blog, nil
}

// Like : 指定ブログにいいねをプラス1
func (b blog) Like(ctx context.Context, title string, isSilent bool) (model.Blog, error) {
	blog, err := b.r.FindByTitle(ctx, title)
	if err != nil {
		if errors.Is(err, repository.ErrBlogRecordNotFound) {
			return model.Blog{}, ErrBlogNotFound
		}
		return model.Blog{}, fmt.Errorf("blogRepository.FindByTitle()内でのエラー: %w", err)
	}

	blog = blog.CountUp()
	blog, err = b.r.Update(ctx, blog)
	if err != nil {
		return model.Blog{}, fmt.Errorf("blogRepository.Update()内でのエラー: %w", err)
	}

	// Slack に通知（リクエスト完了後も継続させるため context の cancel を切り離す）
	if !isSilent {
		notifyCtx := context.WithoutCancel(ctx)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error(notifyCtx, fmt.Errorf("panic in SendLikeNotification: %v", r))
				}
			}()
			if err := b.sg.SendLikeNotification(notifyCtx, blog); err != nil {
				log.Error(notifyCtx, fmt.Errorf("failed to SendLikeNotification: %w", err))
			}
		}()
	}

	return blog, nil
}

package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Blog : Blog用REST Handlerのインターフェース
type Blog interface {
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Like(w http.ResponseWriter, r *http.Request)
}

type blog struct {
	usecase usecase.Blog
}

// NewBlog : Blog用REST Handlerを取得
func NewBlog(u usecase.Blog) Blog {
	return &blog{
		usecase: u,
	}
}

// BlogResponse : Blog用共通レスポンス
type BlogResponse struct {
	Blog *model.Blog `json:"blog,omitempty"`
	errorResponse
}

// Create : ブログ情報を新規作成
func (b blog) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"title" validate:"required,max=50"`
	}

	ctx := r.Context()
	req := request{}
	resp := BlogResponse{}
	if err := bindReqWithValidate(ctx, &req, r); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusBadRequest)
		return
	}

	blog, err := b.usecase.Create(ctx, req.Title)
	if err != nil {
		errInfo := fmt.Errorf("BlogUseCase.Create()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}
	resp.Blog = blog

	DoResponse(w, resp, http.StatusCreated)
}

// Show : ブログ情報を1件取得
func (b blog) Show(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `validate:"required,max=50"`
	}

	ctx := r.Context()
	req := request{}
	resp := BlogResponse{}
	if err := bindReqWithValidate(ctx, &req, nil); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusBadRequest)
		return
	}

	blog, err := b.usecase.Show(ctx, req.Title)
	if err != nil {
		errInfo := fmt.Errorf("BlogUseCase.Show()でエラー: %w", err)
		app.Logger.Println(errInfo)

		if errors.Is(err, usecase.ErrBlogNotFound) {
			DoResponse(w, resp, http.StatusNoContent)
			return
		}

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Blog = blog
	DoResponse(w, resp, http.StatusOK)
}

// Like : 指定ブログにいいねをプラス1
func (b blog) Like(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `validate:"required,max=50"`
	}

	ctx := r.Context()
	req := request{}
	resp := BlogResponse{}
	if err := bindReqWithValidate(ctx, &req, nil); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusBadRequest)
		return
	}

	blog, err := b.usecase.Like(ctx, req.Title)
	if err != nil {
		errInfo := fmt.Errorf("BlogUseCase.Like()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()

		if errors.Is(err, usecase.ErrBlogNotFound) {
			DoResponse(w, resp, http.StatusNoContent)
			return
		}

		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}

	resp.Blog = blog
	DoResponse(w, resp, http.StatusOK)
}

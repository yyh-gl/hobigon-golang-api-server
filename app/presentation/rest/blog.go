package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
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

// Create : ブログ情報を新規作成
func (b blog) Create(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Title string `validate:"required,max=50"`
		}
		response struct {
			Title string `json:"title"`
			Count int    `json:"count"`
		}
	)

	ctx := r.Context()
	req := request{}
	if err := bindReqWithValidate(ctx, r, &req); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate() > %w", err)
		app.Logger.Println(errInfo)

		DoResponse(w, errBadRequest, http.StatusBadRequest)
		return
	}

	blog, err := b.usecase.Create(ctx, req.Title)
	if err != nil {
		// TODO: 全て500エラーにしているのでより詳細なエラーを出す（重複エラーとか）

		errInfo := fmt.Errorf("BlogUseCase.Create()でエラー: %w", err)
		app.Logger.Println(errInfo)

		DoResponse(w, errInterServerError, http.StatusInternalServerError)
		return
	}

	resp := response{
		// 簡略化のためにドメインモデルを直接参照
		Title: blog.Title().String(),
		Count: blog.Count().Int(),
	}
	DoResponse(w, resp, http.StatusCreated)
}

// Show : ブログ情報を1件取得
func (b blog) Show(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Title string `validate:"required,max=50"`
		}
		response struct {
			Title string `json:"title"`
			Count int    `json:"count"`
		}
	)

	ctx := r.Context()

	var req request
	if err := bindReqWithValidate(ctx, mux.Vars(r), &req); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate() > %w", err)
		app.Logger.Println(errInfo)

		DoResponse(w, errBadRequest, http.StatusBadRequest)
		return
	}

	blog, err := b.usecase.Show(ctx, req.Title)
	if err != nil {
		errInfo := fmt.Errorf("BlogUseCase.Show()でエラー: %w", err)
		app.Logger.Println(errInfo)

		if errors.Is(err, usecase.ErrBlogNotFound) {
			DoResponse(w, errNotFound, http.StatusNotFound)
			return
		}

		DoResponse(w, errInterServerError, http.StatusInternalServerError)
		return
	}

	resp := response{
		// 簡略化のためにドメインモデルを直接参照
		Title: blog.Title().String(),
		Count: blog.Count().Int(),
	}
	DoResponse(w, resp, http.StatusOK)
}

// Like : 指定ブログにいいねをプラス1
func (b blog) Like(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Title string `validate:"required,max=50"`
		}
		response struct {
			Title string `json:"title"`
			Count int    `json:"count"`
		}
	)

	ctx := r.Context()

	var req request
	if err := bindReqWithValidate(ctx, mux.Vars(r), &req); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate() > %w", err)
		app.Logger.Println(errInfo)

		DoResponse(w, errBadRequest, http.StatusBadRequest)
		return
	}

	isSilent, _ := strconv.ParseBool(r.Header.Get("x-is-silent"))

	blog, err := b.usecase.Like(ctx, req.Title, isSilent)
	if err != nil {
		errInfo := fmt.Errorf("BlogUseCase.Like()でエラー: %w", err)
		app.Logger.Println(errInfo)

		if errors.Is(err, usecase.ErrBlogNotFound) {
			DoResponse(w, nil, http.StatusNoContent)
			return
		}

		DoResponse(w, errInterServerError, http.StatusInternalServerError)
		return
	}

	resp := response{
		// 簡略化のためにドメインモデルを直接参照
		Title: blog.Title().String(),
		Count: blog.Count().Int(),
	}
	DoResponse(w, resp, http.StatusOK)
}

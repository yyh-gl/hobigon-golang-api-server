package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"

	"github.com/yyh-gl/hobigon-golang-api-server/usecase"

	"github.com/yyh-gl/hobigon-golang-api-server/context"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

// CreateBlogHandler はブログデータを新規に作成
func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	type request struct {
		Title string `json:"title"`
	}

	type blog struct {
		ID        uint       `json:"id,omitempty"`
		Title     string     `json:"title,omitempty"`
		Count     *int       `json:"count,omitempty"`
		CreatedAt *time.Time `json:"created_at,omitempty"`
		UpdatedAt *time.Time `json:"updated_at,omitempty"`
		DeletedAt *time.Time `json:"deleted_at,omitempty"`
	}

	type response struct {
		IsSuccess bool   `json:"is_success"`
		Error     string `json:"error,omitempty"`
		Blog      *blog  `json:"blog,omitempty"`
	}

	res := response{
		IsSuccess: true,
	}

	req, err := decodeRequest(r, request{})
	if err != nil {
		logger.Println(err)

		res.IsSuccess = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer r.Body.Close()

	var b *model.Blog
	if res.IsSuccess {
		params := usecase.CreateBlogParams{
			Title: req["title"].(string),
		}
		b, err = usecase.CreateBlogUseCase(r.Context(), params)
		if err != nil {
			res.IsSuccess = false
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		}

		if b != nil {
			blogRes := blog(*b)
			res.Blog = &blogRes
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Println(err)

		res.IsSuccess = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CreateBlogHandler はブログデータを新規に作成
func GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()
	ps := context.FetchRequestParams(ctx)

	blogRepository := repository.NewBlogRepository()

	blog, err := blogRepository.SelectByTitle(ps.ByName("title"))
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		switch err {
		case gorm.ErrRecordNotFound:
			http.Error(w, "Record Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(blog); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func LikeBlogHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()
	ps := context.FetchRequestParams(ctx)

	blogRepository := repository.NewBlogRepository()
	slackGateway := gateway.NewSlackGateway()

	blog, err := blogRepository.SelectByTitle(ps.ByName("title"))
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		switch err {
		case gorm.ErrRecordNotFound:
			http.Error(w, "Record Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Count をプラス1
	addedCount := *blog.Count + 1
	blog.Count = &addedCount
	blog, err = blogRepository.Update(blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Slack に通知
	err = slackGateway.SendLikeNotify(blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(blog); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

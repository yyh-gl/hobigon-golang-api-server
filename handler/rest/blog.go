package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/context"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"title"`
	}

	logger := app.Logger

	blogRepository := repository.NewBlogRepository()

	req, err := decodeRequest(r, request{})
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	blog := model.Blog{
		Title: req["title"].(string),
	}
	blog, err = blogRepository.Create(blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

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

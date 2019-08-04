package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

type CreateBlogRequest struct {
	Title string `json:"title"`
}

func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value("logger").(log.Logger)

	blogRepository := repository.NewBlogRepository()

	// TODO: デコード処理を共通化
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var createBlogRequest CreateBlogRequest
	err = json.Unmarshal(body, &createBlogRequest)
	if err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var blog model.Blog
	blog.Title = createBlogRequest.Title
	blog, err = blogRepository.Create(ctx, blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value("logger").(log.Logger)

	q := r.URL.Query()
	title := q.Get("title")

	blogRepository := repository.NewBlogRepository()
	blog, err := blogRepository.SelectByTitle(ctx, title)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

type LikeBlogRequest struct {
	Title string `json:"title"`
}

func LikeBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value("logger").(log.Logger)

	blogRepository := repository.NewBlogRepository()

	// TODO: デコード処理を共通化
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var likeBlogRequest LikeBlogRequest
	err = json.Unmarshal(body, &likeBlogRequest)
	if err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	blog, err := blogRepository.SelectByTitle(ctx, likeBlogRequest.Title)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Count をプラス1
	addedCount := *blog.Count + 1
	blog.Count = &addedCount
	blog, err = blogRepository.Update(ctx, blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

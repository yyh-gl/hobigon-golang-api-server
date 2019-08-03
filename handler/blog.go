package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

type CreateBlogRequest struct {
	Title string `json:"title"`
}

type CreateBlogResponse struct {
	Blog model.Blog `json:"blog"`
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

	res := CreateBlogResponse{
		Blog: blog,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

type GetBlogResponse struct {
	Blog model.Blog `json:"blog"`
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

	res := GetBlogResponse{
		Blog: blog,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

type PutBlogRequest struct {
	Title string `json:"title"`
}

type PutBlogResponse struct {
	Blog model.Blog `json:"blog"`
}

func PutBlogHandler(w http.ResponseWriter, r *http.Request) {
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

	var putBlogRequest PutBlogRequest
	err = json.Unmarshal(body, &putBlogRequest)
	if err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	blog, err := blogRepository.SelectByTitle(ctx, putBlogRequest.Title)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	fmt.Println("========================")
	fmt.Println(blog.Title)
	fmt.Println("========================")

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

	res := PutBlogResponse{
		Blog: blog,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

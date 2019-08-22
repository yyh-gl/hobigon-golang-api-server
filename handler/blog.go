package handler

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
	"io/ioutil"
	"log"
	"net/http"
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
	ps := ctx.Value("params").(httprouter.Params)

	blogRepository := repository.NewBlogRepository()

	blog, err := blogRepository.SelectByTitle(ctx, ps.ByName("title"))
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
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

// TODO: いいねされたときの通知機能をつける
func LikeBlogHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value("logger").(log.Logger)
	ps := ctx.Value("params").(httprouter.Params)

	blogRepository := repository.NewBlogRepository()

	blog, err := blogRepository.SelectByTitle(ctx, ps.ByName("title"))
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
	blog, err = blogRepository.Update(ctx, blog)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if err := json.NewEncoder(w).Encode(blog); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

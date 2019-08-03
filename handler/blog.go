package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

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

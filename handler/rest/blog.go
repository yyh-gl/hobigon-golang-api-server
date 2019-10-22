package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/context"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

type blog struct {
	ID        uint       `json:"id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Count     *int       `json:"count,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TODO: OK, Error 部分は共通レスポンスにする
type blogResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Blog  *blog  `json:"blog,omitempty"`
}

// CreateBlogHandler はブログデータを新規に作成
func CreateBlogHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"title"`
	}

	logger := app.Logger

	res := blogResponse{
		OK: true,
	}

	req, err := decodeRequest(r, request{})
	if err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer r.Body.Close()

	var b *model.Blog
	if res.OK {
		b, err = usecase.CreateBlogUseCase(r.Context(), req["title"].(string))
		if err != nil {
			res.OK = false
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetBlogHandler はブログデータを1件取得
func GetBlogHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()
	ps := context.FetchRequestParams(ctx)

	res := blogResponse{
		OK: true,
	}

	b, err := usecase.GetBlogUseCase(ctx, ps.ByName("title"))
	if err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()

		switch err.Error() {
		case "record not found":
			// レコードが存在しないときは空のデータを返す
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if b != nil {
		blogRes := blog(*b)
		res.Blog = &blogRes
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// LikeBlogHandler は指定ブログにいいねをプラス1
func LikeBlogHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()
	ps := context.FetchRequestParams(ctx)

	res := blogResponse{
		OK: true,
	}

	b, err := usecase.LikeBlogUseCase(ctx, ps.ByName("title"))
	if err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()

		switch err.Error() {
		case "record not found":
			// レコードが存在しないときは空のデータを返す
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if b != nil {
		blogRes := blog(*b)
		res.Blog = &blogRes
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

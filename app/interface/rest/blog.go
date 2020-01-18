package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	u usecase.Blog
}

// NewBlog : Blog用REST Handlerを取得
func NewBlog(u usecase.Blog) Blog {
	return &blog{
		u: u,
	}
}

// blogResponse : Blog用共通レスポンス
// TODO: OK, Error 部分は共通レスポンスにする
type blogResponse struct {
	OK    bool       `json:"ok"`
	Error string     `json:"error,omitempty"`
	Blog  model.Blog `json:"blog,omitempty"`
}

// Create : ブログ情報を新規作成
func (b blog) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"title"`
	}

	res := blogResponse{
		OK: true,
	}

	req, err := decodeRequest(r, request{})
	if err != nil {
		app.Logger.Println(fmt.Errorf("decodeRequest()でエラー: %w", err))

		res.OK = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer func() { _ = r.Body.Close() }()

	var createdBlog *model.Blog
	if res.OK {
		createdBlog, err = b.u.Create(r.Context(), req["title"].(string))
		if err != nil {
			app.Logger.Println(fmt.Errorf("BlogUseCase.Create()でエラー: %w", err))

			res.OK = false
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			res.Blog = *createdBlog
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		app.Logger.Println(fmt.Errorf("json.NewEncoder().Encode()でエラー: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Show : ブログ情報を1件取得
func (b blog) Show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO: httprouterに依存することについて考える
	ps := httprouter.ParamsFromContext(ctx)

	res := blogResponse{
		OK: true,
	}

	blog, err := b.u.Show(ctx, ps.ByName("title"))
	if err != nil {
		app.Logger.Println(err)

		res.OK = false
		res.Error = err.Error()

		switch err.Error() {
		case "record not found":
			// レコードが存在しないときは空の情報を返す
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		res.Blog = *blog
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Like : 指定ブログにいいねをプラス1
func (b blog) Like(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ps := httprouter.ParamsFromContext(ctx)

	res := blogResponse{
		OK: true,
	}

	blog, err := b.u.Like(ctx, ps.ByName("title"))
	if err != nil {
		app.Logger.Println(err)

		res.OK = false
		res.Error = err.Error()

		switch err.Error() {
		case "record not found":
			// レコードが存在しないときは空の情報を返す
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		res.Blog = *blog
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

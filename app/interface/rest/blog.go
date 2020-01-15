package rest

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

//////////////////////////////////////////////////
// NewBlogHandler
//////////////////////////////////////////////////

// BlogHandler : ブログ用のハンドラーインターフェース
type BlogHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Like(w http.ResponseWriter, r *http.Request)
}

type blogHandler struct {
	bu usecase.BlogUseCase
}

// NewBlogHandler : ブログ用のハンドラーを取得
func NewBlogHandler(bu usecase.BlogUseCase) BlogHandler {
	return &blogHandler{
		bu: bu,
	}
}

// TODO: OK, Error 部分は共通レスポンスにする
type blogResponse struct {
	OK    bool          `json:"ok"`
	Error string        `json:"error,omitempty"`
	Blog  blog.BlogJSON `json:"blog,omitempty"`
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : ブログ情報を新規作成
func (bh blogHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"title"`
	}

	res := blogResponse{
		OK: true,
	}

	req, err := decodeRequest(r, request{})
	if err != nil {
		app.Logger.Println(err)

		res.OK = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer func() { _ = r.Body.Close() }()

	var b *blog.Blog
	if res.OK {
		b, err = bh.bu.Create(r.Context(), req["title"].(string))
		if err != nil {
			app.Logger.Println(err)
			res.OK = false
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// JSON 形式に変換
			res.Blog = b.JSONSerialize()
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		app.Logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

//////////////////////////////////////////////////
// Show
//////////////////////////////////////////////////

// Show : ブログ情報を1件取得
func (bh blogHandler) Show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO: httprouterに依存することについて考える
	ps := httprouter.ParamsFromContext(ctx)

	res := blogResponse{
		OK: true,
	}

	b, err := bh.bu.Show(ctx, ps.ByName("title"))
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
		// JSON 形式に変換
		res.Blog = b.JSONSerialize()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

//////////////////////////////////////////////////
// Like
//////////////////////////////////////////////////

// Like : 指定ブログにいいねをプラス1
func (bh blogHandler) Like(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ps := httprouter.ParamsFromContext(ctx)

	res := blogResponse{
		OK: true,
	}

	b, err := bh.bu.Like(ctx, ps.ByName("title"))
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
		// JSON 形式に変換
		res.Blog = b.JSONSerialize()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

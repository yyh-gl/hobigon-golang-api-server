package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/context"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
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

// response : ブログ用共通正常時レスポンス
type response struct {
	Blog blog.BlogJSON `json:"blog"`
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : ブログ情報を新規作成
func (bh blogHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Title string `json:"title"`
	}

	logger := app.Logger

	errRes := new(errorResponse)
	req, err := decodeRequest(r, request{})
	if err != nil {
		logger.Println(err)
		errRes.Detail = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer func() { _ = r.Body.Close() }()

	res := new(response)
	if errRes.Detail == "" {
		b := new(blog.Blog)
		b, err = bh.bu.Create(r.Context(), req["title"].(string))
		if err != nil {
			errRes.Detail = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// JSON 形式に変換
			res.Blog = b.JSONSerialize()
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

//////////////////////////////////////////////////
// Show
//////////////////////////////////////////////////

// Show : ブログ情報を1件取得
func (bh blogHandler) Show(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()
	ps := context.FetchRequestParams(ctx)

	res := new(response)
	errRes := new(errorResponse)
	b, err := bh.bu.Show(ctx, ps.ByName("title"))
	if err != nil {
		logger.Println(err)

		errRes.Detail = err.Error()
		if errRes.Detail == usecase.ErrRecordNotFound {
			// レコードが存在しないときは空の情報を返す
			w.WriteHeader(http.StatusNotFound)
		}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// JSON 形式に変換
		res.Blog = b.JSONSerialize()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

//////////////////////////////////////////////////
// Like
//////////////////////////////////////////////////

// Like : 指定ブログにいいねをプラス1
func (bh blogHandler) Like(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()
	ps := context.FetchRequestParams(ctx)

	res := new(response)
	errRes := new(errorResponse)
	b, err := bh.bu.Like(ctx, ps.ByName("title"))
	if err != nil {
		logger.Println(err)

		errRes.Detail = err.Error()
		if errRes.Detail == usecase.ErrRecordNotFound {
			// レコードが存在しないときは空の情報を返す
			w.WriteHeader(http.StatusNotFound)
		}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// JSON 形式に変換
		res.Blog = b.JSONSerialize()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

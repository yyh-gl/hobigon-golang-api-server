package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

func TestBlogHandler_Create(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}

	tests := []struct {
		title string
		want  want
	}{
		{ // 正常系
			title: "sample-blog-title",
			want: want{
				body:       `{"title":"sample-blog-title","count":0}`,
				statusCode: http.StatusCreated,
			},
		},
		{ // 正常系：50文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
			want: want{
				body:       `{"title":"hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title","count":0}`,
				statusCode: http.StatusCreated,
			},
		},
		{ // 異常系：51文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title-over",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：タイトルを渡さない
			title: "",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：duplicate
			title: "duplicate-blog-title",
			want: want{
				body:       `{"error":{"detail":"サーバ内でエラーが発生しました"}}`,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	// 重複データ登録時に使用するテストデータを追加
	test.CreateBlog(DIContainer.DB, "duplicate-blog-title")

	for _, tt := range tests {
		reqBody := strings.NewReader(`{"title":"` + tt.title + `"}`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/blogs", reqBody)
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)

		if c := rr.Code; c != tt.want.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				c, tt.want.statusCode)
		}

		respBody := strings.TrimRight(rr.Body.String(), "\n")
		if respBody != tt.want.body {
			t.Errorf("handler returned unexpected body: got %v want %v",
				respBody, tt.want.body)
		}
	}
}

func TestBlogHandler_Show(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}

	tests := []struct {
		title string
		want  want
	}{
		{ // 正常系
			title: "sample-blog-title",
			want: want{
				body:       `{"title":"sample-blog-title","count":0}`,
				statusCode: http.StatusOK,
			},
		},
		{ // 正常系：存在しないブログ
			title: "sample-blog-title2",
			want: want{
				body:       "null",
				statusCode: http.StatusNotFound,
			},
		},
		{ // 正常系：50文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
			want: want{
				body:       `{"title":"hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title","count":0}`,
				statusCode: http.StatusOK,
			},
		},
		{ // 異常系：51文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title-over",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	// テストデータを追加
	test.CreateBlog(DIContainer.DB, "sample-blog-title")
	test.CreateBlog(DIContainer.DB, "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title")

	for _, tt := range tests {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/blogs/"+tt.title, nil)
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)

		if c := rr.Code; c != tt.want.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				c, tt.want.statusCode)
		}

		b := strings.TrimRight(rr.Body.String(), "\n")
		if b != tt.want.body {
			t.Errorf("handler returned unexpected body: got %v want %v",
				b, tt.want.body)
		}
	}
}

func TestBlogHandler_Like(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}

	tests := []struct {
		title string
		want  want
	}{
		{ // 正常系
			title: "sample-blog-title",
			want: want{
				body:       `{"title":"sample-blog-title","count":1}`,
				statusCode: http.StatusOK,
			},
		},
		{ // 正常系：存在しないブログ
			title: "sample-blog-title2",
			want: want{
				body:       "null",
				statusCode: http.StatusNoContent,
			},
		},
		{ // 正常系：50文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
			want: want{
				body:       `{"title":"hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title","count":1}`,
				statusCode: http.StatusOK,
			},
		},
		{ // 異常系：51文字タイトル
			// TODO: ドメインモデル用テストを作って、NewTitle()におけるバリデーションが動作するか確認
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title-over",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
	}

	// テストデータを追加
	test.CreateBlog(DIContainer.DB, "sample-blog-title")
	test.CreateBlog(DIContainer.DB, "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title")

	for _, tt := range tests {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/blogs/"+tt.title+"/like", nil)
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)

		if c := rr.Code; c != tt.want.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				c, tt.want.statusCode)
		}

		b := strings.TrimRight(rr.Body.String(), "\n")
		if b != tt.want.body {
			t.Errorf("handler returned unexpected body: got %v want %v",
				b, tt.want.body)
		}
	}
}

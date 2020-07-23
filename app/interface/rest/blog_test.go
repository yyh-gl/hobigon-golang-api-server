package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

//func createBlog(c *test.Client, title string) {
//	c.AddRoute(http.MethodPost, "/api/v1/blogs", c.DIContainer.HandlerBlog.Create)
//
//	body := `{
//  "title": "` + title + `"
//}`
//	_ = c.Post("/api/v1/blogs", body)
//}

//func TestBlogHandler_Create(t *testing.T) {
//	c := test.NewClient()
//	defer func() { _ = c.DIContainer.DB.Close() }()
//
//	c.AddRoute(http.MethodPost, "/api/v1/blogs", c.DIContainer.HandlerBlog.Create)
//
//	testCases := []struct {
//		title string
//		want  string
//		err   string
//	}{
//		{ // 正常系
//			title: "sample-blog-title",
//			want:  "sample-blog-title",
//			err:   "",
//		},
//		{ // 正常系：50文字タイトル
//			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
//			want:  "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
//			err:   "",
//		},
//		{ // 異常系：51文字タイトル
//			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title-over",
//			want:  "",
//			err:   "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Title' Error:Field validation for 'Title' failed on the 'max' tag",
//		},
//		{ // 異常系：タイトルを渡さない
//			title: "",
//			want:  "",
//			err:   "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Title' Error:Field validation for 'Title' failed on the 'required' tag",
//		},
//		{ // 異常系：duplicate
//			title: "duplicate-blog-title",
//			want:  "",
//			err:   "BlogUseCase.Create()でエラー: blogRepository.Create()内でのエラー: gorm.Create(blog)内でのエラー: UNIQUE constraint failed: blog_posts.title",
//		},
//	}
//
//	// 重複データ登録時に使用するテストデータを追加
//	createBlog(c, "duplicate-blog-title")
//
//	for _, tc := range testCases {
//		body := `{
// "title": "` + tc.title + `"
//}`
//		rec := c.Post("/api/v1/blogs", body)
//		resp := rest.BlogResponse{}
//		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
//
//		if tc.err == "" {
//			assert.Equal(t, tc.want, resp.Blog.Title().String())
//			assert.Equal(t, "", resp.Error)
//		} else {
//			assert.Equal(t, (*blog.Blog)(nil), resp.Blog)
//			assert.Equal(t, tc.err, resp.Error)
//		}
//	}
//}

func TestBlogHandler_Show(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}

	tests := []struct {
		title string
		want  want
		err   string
	}{
		{ // 正常系
			title: "sample-blog-title",
			want: want{
				body:       `{"Title":"sample-blog-title","Count":0}`,
				statusCode: http.StatusOK,
			},
			err: "",
		},
		{ // 正常系：存在しないブログ
			title: "sample-blog-title2",
			want: want{
				body:       "null",
				statusCode: http.StatusNotFound,
			},
			err: "",
		},
		{ // 正常系：50文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
			want: want{
				body:       `{"Title":"hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title","Count":0}`,
				statusCode: http.StatusOK,
			},
			err: "",
		},
		{ // 異常系：51文字タイトル
			title: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title-over",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
			err: "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Title' Error:Field validation for 'Title' failed on the 'max' tag",
		},
	}

	// テストデータを追加
	test.CreateBlog(DIContainer.DB, "sample-blog-title")
	test.CreateBlog(DIContainer.DB, "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title")
	defer func() {
		test.DeleteBlog(DIContainer.DB, "sample-blog-title")
		test.DeleteBlog(DIContainer.DB, "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title")
	}()

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

//func TestBlogHandler_Like(t *testing.T) {
//	c := test.NewClient()
//	defer func() { _ = c.DIContainer.DB.Close() }()
//
//	c.AddRoute(http.MethodPost, "/api/v1/blogs/:title/like", c.DIContainer.HandlerBlog.Like)
//
//	testCases := []struct {
//		title     string
//		wantTitle string
//		wantCount int
//		err       string
//	}{
//		{ // 正常系
//			title:     "sample-blog-title",
//			wantTitle: "sample-blog-title",
//			wantCount: 1,
//			err:       "",
//		},
//		{ // 正常系：存在しないブログ
//			title:     "sample-blog-title2",
//			wantTitle: "",
//			wantCount: 0,
//			err:       "",
//		},
//		{ // 正常系：50文字タイトル
//			title:     "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
//			wantTitle: "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title",
//			wantCount: 1,
//			err:       "",
//		},
//		{ // 異常系：51文字タイトル
//			title:     "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title-over",
//			wantTitle: "",
//			wantCount: 0,
//			// TODO: ドメインモデル用テストを作って、NewTitle()におけるバリデーションが動作するか確認
//			err: "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Title' Error:Field validation for 'Title' failed on the 'max' tag",
//		},
//	}
//
//	// テストデータを追加
//	createBlog(c, "sample-blog-title")
//	createBlog(c, "hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-hoge-title")
//
//	for _, tc := range testCases {
//		rec := c.Post("/api/v1/blogs/"+tc.title+"/like", "")
//		resp := rest.BlogResponse{}
//		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
//
//		if tc.err == "" {
//			if tc.wantTitle != "" {
//				assert.Equal(t, tc.wantTitle, resp.Blog.Title().String())
//			}
//			if tc.wantCount != 0 {
//				assert.Equal(t, tc.wantCount, resp.Blog.Count().Int())
//			}
//			assert.Equal(t, "", resp.Error)
//		} else {
//			assert.Equal(t, (*blog.Blog)(nil), resp.Blog)
//			assert.Equal(t, tc.err, resp.Error)
//		}
//	}
//}

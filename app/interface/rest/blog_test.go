package rest_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

type blogResponse struct {
	OK    bool      `json:"ok"`
	Error string    `json:"error,omitempty"`
	Blog  blog.Blog `json:"blog,omitempty"`
}

func createBlog(c *test.Client, title string) {
	c.AddRoute(http.MethodPost, "/api/v1/blogs", c.DIContainer.HandlerBlog.Create)

	body := `{
  "title": "` + title + `"
}`
	_ = c.Post("/api/v1/blogs", body)
}

func TestBlogHandler_Create(t *testing.T) {
	c := test.NewClient()
	defer func() { _ = c.DIContainer.DB.Close() }()

	c.AddRoute(http.MethodPost, "/api/v1/blogs", c.DIContainer.HandlerBlog.Create)

	testCases := []struct {
		title string
		want  string
		err   string
	}{
		{ // 正常系
			title: "sample-blog-title",
			want:  "sample-blog-title",
			err:   "",
		},
		{ // 異常系：タイトルを渡さない
			title: "",
			want:  "",
			err:   "バリデーションエラー：タイトルは必須です",
		},
		{ // 異常系：duplicate
			title: "duplicate-blog-title",
			want:  "",
			err:   "blogRepository.Create()内でのエラー: gorm.Create(blog)内でのエラー: UNIQUE constraint failed: blog_posts.title",
		},
	}

	// 重複データ登録時に使用するテストデータを追加
	createBlog(c, "duplicate-blog-title")

	for _, tc := range testCases {
		body := `{
 "title": "` + tc.title + `"
}`
		rec := c.Post("/api/v1/blogs", body)
		resp := new(blogResponse)
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)

		assert.Equal(t, tc.want, resp.Blog.Title())
		assert.Equal(t, tc.err, resp.Error)
	}
}

func TestBlogHandler_Show(t *testing.T) {
	c := test.NewClient()
	defer func() { _ = c.DIContainer.DB.Close() }()

	c.AddRoute(http.MethodGet, "/api/v1/blogs/:title", c.DIContainer.HandlerBlog.Show)

	testCases := []struct {
		title string
		want  string
		err   string
	}{
		{ // 正常系
			title: "sample-blog-title",
			want:  "sample-blog-title",
			err:   "",
		},
		{ // 正常系：存在しないブログ
			title: "sample-blog-title2",
			want:  "",
			err:   "blogRepository.SelectByTitle()内でのエラー: gorm.First(blog)内でのエラー: record not found",
		},
		{ // 異常系：タイトルを渡さない
			title: "",
			want:  "",
			err:   "",
		},
	}

	// テストデータを追加
	createBlog(c, "sample-blog-title")

	for _, tc := range testCases {
		rec := c.Get("/api/v1/blogs/" + tc.title)
		resp := new(blogResponse)
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)

		assert.Equal(t, tc.want, resp.Blog.Title())
		assert.Equal(t, tc.err, resp.Error)
	}
}

func TestBlogHandler_Like(t *testing.T) {
	c := test.NewClient()
	defer func() { _ = c.DIContainer.DB.Close() }()

	c.AddRoute(http.MethodPost, "/api/v1/blogs/:title/like", c.DIContainer.HandlerBlog.Like)

	testCases := []struct {
		title     string
		wantTitle string
		wantCount int
		err       string
	}{
		{ // 正常系
			title:     "sample-blog-title",
			wantTitle: "sample-blog-title",
			wantCount: 1,
			err:       "",
		},
		{ // 正常系：存在しないブログ
			title:     "sample-blog-title2",
			wantTitle: "",
			wantCount: 0,
			err:       "blogRepository.SelectByTitle()内でのエラー: gorm.First(blog)内でのエラー: record not found",
		},
		{ // 異常系：タイトルを渡さない
			title:     "",
			wantTitle: "",
			wantCount: 0,
			err:       "blogRepository.SelectByTitle()内でのエラー: gorm.First(blog)内でのエラー: record not found",
		},
	}

	// テストデータを追加
	createBlog(c, "sample-blog-title")

	for _, tc := range testCases {
		rec := c.Post("/api/v1/blogs/"+tc.title+"/like", "")
		resp := new(blogResponse)
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)

		assert.Equal(t, tc.wantTitle, resp.Blog.Title())
		if tc.wantCount != 0 {
			assert.Equal(t, tc.wantCount, *resp.Blog.Count())
		}
		assert.Equal(t, tc.err, resp.Error)
	}
}

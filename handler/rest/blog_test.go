package rest_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/blog"
)

func TestBlogHandler_Create(t *testing.T) {
}

func TestBlogHandler_Show(t *testing.T) {
	type blogResponse struct {
		OK    bool          `json:"ok"`
		Error string        `json:"error,omitempty"`
		Blog  blog.BlogJSON `json:"blog,omitempty"`
	}

	testCases := []struct {
		title string
		want  string
		err   error
	}{
		{ // 正常系
			title: "sample-blog-title",
			want:  "sample-blog-title",
			err:   nil,
		},
		{ // 正常系：存在しないタイトル
			title: "no-exist-sample-blog-title",
			want:  "",
			// TODO: err を返さないようにする
			err: errors.New("blogRepository.SelectByTitle()内でのエラー: gorm.First(blog)内でのエラー: record not found"),
		},
	}
	for _, tc := range testCases {
		url := fmt.Sprintf("http://localhost:3000/api/v1/blogs/%s", tc.title)
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)

		res := new(blogResponse)
		_ = json.Unmarshal(byteArray, &res)

		assert.Equal(t, tc.want, res.Blog.Title)
		if res.Error != "" {
			assert.Equal(t, tc.err.Error(), res.Error)
		}
	}
}

func TestBlogHandler_Like(t *testing.T) {

}

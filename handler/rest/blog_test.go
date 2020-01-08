package rest_test

import (
	"errors"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

type blogResponse struct {
	OK    bool          `json:"ok"`
	Error string        `json:"error,omitempty"`
	Blog  blog.BlogJSON `json:"blog,omitempty"`
}

func TestBlogHandler_Create(t *testing.T) {
	test.NewTestDBConnect()
	test.ResetTable(test.BlogTable)

	type params struct {
		Title string `json:"title"`
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
		{ // 異常系：タイトルを渡さない
			title: "",
			want:  "",
			err:   errors.New("unexpected end of JSON input"),
		},
	}

	url := "http://localhost:3000/api/v1/blogs"
	for _, tc := range testCases {
		p := params{Title: tc.title}
		resp := test.Post(url, p)

		//assert.Equal(t, tc.want, res.Blog.Title)
		//if res.Error != "" {
		//	assert.Equal(t, tc.err.Error(), res.Error)
		//}
	}
}

func TestBlogHandler_Show(t *testing.T) {
	//test.NewTestDBConnect()
	//test.ResetTable(test.BlogTable)
	//
	//testCases := []struct {
	//	title string
	//	want  string
	//	err   error
	//}{
	//	{ // 正常系
	//		title: "sample-blog-title",
	//		want:  "sample-blog-title",
	//		err:   nil,
	//	},
	//	{ // 正常系：存在しないタイトル
	//		title: "no-exist-sample-blog-title",
	//		want:  "",
	//		// TODO: err を返さないようにする
	//		err: errors.New("blogRepository.SelectByTitle()内でのエラー: gorm.First(blog)内でのエラー: record not found"),
	//	},
	//}
	//for _, tc := range testCases {
	//	url := fmt.Sprintf("http://localhost:3000/api/v1/blogs/%s", tc.title)
	//	resp, _ := http.Get(url)
	//	defer resp.Body.Close()
	//	byteArray, _ := ioutil.ReadAll(resp.Body)
	//
	//	res := new(blogResponse)
	//	_ = json.Unmarshal(byteArray, &res)
	//
	//	assert.Equal(t, tc.want, res.Blog.Title)
	//	if res.Error != "" {
	//		assert.Equal(t, tc.err.Error(), res.Error)
	//	}
	//}
}

func TestBlogHandler_Like(t *testing.T) {
	//testCases := []struct {
	//	title string
	//	want  string
	//	err   error
	//}{
	//	{ // 正常系
	//		title: "sample-blog-title",
	//		want:  "sample-blog-title",
	//		err:   nil,
	//	},
	//	{ // 正常系：存在しないタイトル
	//		title: "no-exist-sample-blog-title",
	//		want:  "",
	//		// TODO: record not found にする
	//		err: errors.New("blogRepository.SelectByTitle()内でのエラー: gorm.First(blog)内でのエラー: record not found"),
	//	},
	//}
	//for _, tc := range testCases {
	//	url := fmt.Sprintf("http://localhost:3000/api/v1/blogs/%s", tc.title)
	//	resp, _ := http.Get(url)
	//	defer resp.Body.Close()
	//	byteArray, _ := ioutil.ReadAll(resp.Body)
	//
	//	before := new(blogResponse)
	//	_ = json.Unmarshal(byteArray, &before)
	//
	//	url = fmt.Sprintf("http://localhost:3000/api/v1/blogs/%s/like", tc.title)
	//	resp, _ = http.Get(url)
	//	defer resp.Body.Close()
	//	byteArray, _ = ioutil.ReadAll(resp.Body)
	//
	//	after := new(blogResponse)
	//	_ = json.Unmarshal(byteArray, &after)
	//
	//	assert.Equal(t, before.Blog.Count, *after.Blog.Count+1)
	//	if after.Error != "" {
	//		assert.Equal(t, tc.err.Error(), after.Error)
	//	}
	//}
}

package blog

import (
	"encoding/json"
	"fmt"
)

// Blog : ブログを表すドメインモデル
type Blog struct {
	f fields
}

type fields struct {
	Title Title `json:"title"`
	Count Count `json:"count"`
}

// NewBlog : Blog ドメインモデルを生成
func NewBlog(title string) (*Blog, error) {
	// Title を生成
	t, err := NewTitle(title)
	if err != nil {
		return nil, fmt.Errorf("NewTitle()内でエラー: %w", err)
	}

	// Count を生成
	c, err := NewCount()
	if err != nil {
		return nil, fmt.Errorf("NewTitle()内でエラー: %w", err)
	}

	return &Blog{
		f: fields{
			Title: *t,
			Count: *c,
		},
	}, nil
}

// NewBlogWithFullParams : パラメータ全指定で Blog ドメインモデルを生成
func NewBlogWithFullParams(title string, count int) (*Blog, error) {
	// Title を生成
	t, err := NewTitle(title)
	if err != nil {
		return nil, fmt.Errorf("NewTitle()内でエラー: %w", err)
	}

	// Count を生成
	c, err := NewCountWithArg(count)
	if err != nil {
		return nil, fmt.Errorf("NewTitle()内でエラー: %w", err)
	}

	return &Blog{
		f: fields{
			Title: *t,
			Count: *c,
		},
	}, nil
}

// Title : title のゲッター
func (b Blog) Title() Title {
	return b.f.Title
}

// Count : count のゲッター
func (b Blog) Count() Count {
	return b.f.Count
}

// CountUp : いいね数をプラス1
func (b Blog) CountUp() *Blog {
	b.f.Count += 1
	return &b
}

// CreateLikeMessage : いいね受信メッセージを生成
func (b Blog) CreateLikeMessage() string {
	return "【" + b.f.Title.String() + "】いいね！（Total: " + b.f.Count.String() + "）"
}

// MarshalJSON : Marshal用関数
// FIXME: ドメインモデル内に持ちたくないが、フィールドを公開したくもないので一旦これでいく。よりよい方法を探す
func (b Blog) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.f)
}

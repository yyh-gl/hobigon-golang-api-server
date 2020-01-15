package blog

import (
	"strconv"
	"time"
)

//////////////////////////////////////////////////
// Blog
//////////////////////////////////////////////////

// Blog : ブログ用のドメインモデル
type Blog struct {
	id        uint
	title     string
	count     *int
	createdAt *time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

// NewBlog : Blog ドメインモデルを生成
func NewBlog(title string) *Blog {
	initCount := 0
	return &Blog{
		title: title,
		count: &initCount,
	}
}

// NewBlogWithFullParams : パラメータ全指定で Blog ドメインモデルを生成
func NewBlogWithFullParams(
	id uint,
	title string,
	count *int,
	createdAt *time.Time,
	updatedAt *time.Time,
	deletedAt *time.Time,
) *Blog {
	return &Blog{
		id:        id,
		title:     title,
		count:     count,
		createdAt: createdAt,
		updatedAt: updatedAt,
		deletedAt: deletedAt,
	}
}

// ID : id のゲッター
func (b Blog) ID() uint {
	return b.id
}

// Title : title のゲッター
func (b Blog) Title() string {
	return b.title
}

// Count : count のゲッター
func (b Blog) Count() *int {
	return b.count
}

// CountUp : いいね数をプラス1
func (b *Blog) CountUp() {
	count := *b.count + 1
	b.count = &count
}

// CreateLikeMessage : いいね受信メッセージを生成
func (b Blog) CreateLikeMessage() string {
	return "【" + b.title + "】いいね！（Total: " + strconv.Itoa(*b.count) + "）"
}

//////////////////////////////////////////////////
// BlogJSON
//////////////////////////////////////////////////

type blogJSONFields struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Count     *int       `json:"count"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// BlogJSON : ブログ用の JSON レスポンス形式の定義
type BlogJSON struct {
	blogJSONFields
}

// JSONSerialize : JSON タグを含む構造体を返却
func (b Blog) JSONSerialize() BlogJSON {
	return BlogJSON{blogJSONFields{
		ID:        b.id,
		Title:     b.title,
		Count:     b.count,
		CreatedAt: b.createdAt,
		UpdatedAt: b.updatedAt,
		DeletedAt: b.deletedAt,
	}}
}

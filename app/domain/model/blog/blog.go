package blog

import (
	"strconv"
	"time"
)

// Blog : ブログを表すドメインモデル
type Blog struct {
	fields
}

type fields struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"` // TODO: VOにする
	Count     *int       `json:"count"` // TODO: ポインタをやめて、VOでデフォルト値を入れるようにする
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewBlog : Blog ドメインモデルを生成
func NewBlog(title string) *Blog {
	initCount := 0
	return &Blog{
		fields{
			Title: title,
			Count: &initCount,
		},
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
		fields{
			ID:        id,
			Title:     title,
			Count:     count,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		},
	}
}

// ID : id のゲッター
func (b Blog) ID() uint {
	return b.fields.ID
}

// Title : title のゲッター
func (b Blog) Title() string {
	return b.fields.Title
}

// Count : count のゲッター
func (b Blog) Count() *int {
	return b.fields.Count
}

// CountUp : いいね数をプラス1
func (b *Blog) CountUp() {
	count := *b.fields.Count + 1
	b.fields.Count = &count
}

// CreateLikeMessage : いいね受信メッセージを生成
func (b Blog) CreateLikeMessage() string {
	return "【" + b.fields.Title + "】いいね！（Total: " + strconv.Itoa(*b.fields.Count) + "）"
}

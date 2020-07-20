package blog

import (
	"encoding/json"
	"fmt"
)

// Blog : ブログを表すドメインモデル
type Blog struct {
	title Title
	count Count
}

// NewBlog : Blog ドメインモデルを生成
func NewBlog(title string) (*Blog, error) {
	t, err := newTitle(title)
	if err != nil {
		return nil, fmt.Errorf("NewTitle()内でエラー: %w", err)
	}

	c, err := newCount()
	if err != nil {
		return nil, fmt.Errorf("NewTitle()内でエラー: %w", err)
	}

	return &Blog{
		title: *t,
		count: *c,
	}, nil
}

// Title : titleのゲッター
func (b Blog) Title() Title {
	return b.title
}

// Count : countのゲッター
func (b Blog) Count() Count {
	return b.count
}

// CountUp : いいね数をプラス1
func (b Blog) CountUp() *Blog {
	b.count += 1
	return &b
}

// CreateLikeMessage : いいね受信メッセージを生成
func (b Blog) CreateLikeMessage() string {
	return "【" + b.title.String() + "】いいね！（Total: " + b.count.String() + "）"
}

// MarshalJSON : Marshal用関数
// FIXME: ドメインモデル内に持ちたくないが、フィールドを公開したくもないので一旦これでいく。よりよい方法を探す
func (b Blog) MarshalJSON() ([]byte, error) {
	return json.Marshal(b)
}

// UnmarshalJSON : Unmarshal用関数
// FIXME: ドメインモデル内に持ちたくないが、フィールドを公開したくもないので一旦これでいく。よりよい方法を探す
//        テストのためにだけに用意しているので、いっそう見直したい
func (b *Blog) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &b)
	if err != nil {
		return fmt.Errorf("Blog.UnmarshalJSON()内でエラー: %w", err)
	}
	return nil
}

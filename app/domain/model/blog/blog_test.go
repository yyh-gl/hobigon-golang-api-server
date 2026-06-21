package blog_test

import (
	"context"
	"strings"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

func TestNewBlog(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"正常系：1文字", "a", false},
		{"正常系：50文字", strings.Repeat("a", 50), false},
		{"異常系：51文字", strings.Repeat("a", 51), true},
		{"正常系：空文字", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := blog.NewBlog(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBlog(%q) error = %v, wantErr %v", tt.title, err, tt.wantErr)
			}
		})
	}
}

func TestTitle_String(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"hello", "hello"},
		{"", ""},
	}
	for _, tt := range tests {
		b, err := blog.NewBlog(tt.title)
		if err != nil {
			t.Fatalf("NewBlog(%q) unexpected error: %v", tt.title, err)
		}
		if got := b.Title().String(); got != tt.want {
			t.Errorf("Title.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestTitle_IsNull(t *testing.T) {
	tests := []struct {
		title string
		want  bool
	}{
		{"", true},
		{"hello", false},
	}
	for _, tt := range tests {
		b, err := blog.NewBlog(tt.title)
		if err != nil {
			t.Fatalf("NewBlog(%q) unexpected error: %v", tt.title, err)
		}
		if got := b.Title().IsNull(); got != tt.want {
			t.Errorf("Title.IsNull() = %v, want %v", got, tt.want)
		}
	}
}

func TestCount_IntAndString(t *testing.T) {
	b, _ := blog.NewBlog("test")
	if got := b.Count().Int(); got != 0 {
		t.Errorf("Count.Int() = %d, want 0", got)
	}
	if got := b.Count().String(); got != "0" {
		t.Errorf("Count.String() = %q, want \"0\"", got)
	}
}

func TestBlog_CountUp(t *testing.T) {
	b, _ := blog.NewBlog("test")
	b = b.CountUp()
	if got := b.Count().Int(); got != 1 {
		t.Errorf("after 1 CountUp: Count = %d, want 1", got)
	}
	b = b.CountUp()
	if got := b.Count().Int(); got != 2 {
		t.Errorf("after 2 CountUp: Count = %d, want 2", got)
	}
}

func TestBlog_CreateLikeMessage(t *testing.T) {
	b, _ := blog.NewBlog("My Blog")
	b = b.CountUp()
	got := b.CreateLikeMessage()
	want := "【My Blog】いいね！（Total: 1）"
	if got != want {
		t.Errorf("CreateLikeMessage() = %q, want %q", got, want)
	}
}

func TestConvertToEntity(t *testing.T) {
	d := dto.BlogDTO{Title: "dto-title", Count: 5}
	b := blog.ConvertToEntity(context.Background(), d)
	if got := b.Title().String(); got != "dto-title" {
		t.Errorf("Title = %q, want \"dto-title\"", got)
	}
	if got := b.Count().Int(); got != 5 {
		t.Errorf("Count = %d, want 5", got)
	}
}

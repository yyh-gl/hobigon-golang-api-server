package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
)

// NotionTaskDTO : Notionで管理しているタスク用のDTO
// FIXME: JSON-to-Go で適当に作った構造体なので必要に応じて修正
type NotionTaskDTO struct {
	Object  string `json:"object"`
	Results []struct {
		Object         string    `json:"object"`
		ID             string    `json:"id"`
		CreatedTime    time.Time `json:"created_time"`
		LastEditedTime time.Time `json:"last_edited_time"`
		CreatedBy      struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		Cover any `json:"cover"`
		Icon  struct {
			Type  string `json:"type"`
			Emoji string `json:"emoji"`
		} `json:"icon"`
		Parent struct {
			Type       string `json:"type"`
			DatabaseID string `json:"database_id"`
		} `json:"parent"`
		Archived   bool `json:"archived"`
		Properties struct {
			DateCreated struct {
				ID          string    `json:"id"`
				Type        string    `json:"type"`
				CreatedTime time.Time `json:"created_time"`
			} `json:"Date Created"`
			Label struct {
				ID          string `json:"id"`
				Type        string `json:"type"`
				MultiSelect []any  `json:"multi_select"`
			} `json:"Label"`
			Deadline struct {
				ID   string `json:"id"`
				Type string `json:"type"`
				Date struct {
					Start    string `json:"start"`
					End      any    `json:"end"`
					TimeZone any    `json:"time_zone"`
				} `json:"date"`
			} `json:"Deadline"`
			Status struct {
				ID     string `json:"id"`
				Type   string `json:"type"`
				Select struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Color string `json:"color"`
				} `json:"select"`
			} `json:"Status"`
			Name struct {
				ID    string `json:"id"`
				Type  string `json:"type"`
				Title []struct {
					Type string `json:"type"`
					Text struct {
						Content string `json:"content"`
						Link    any    `json:"link"`
					} `json:"text"`
					Annotations struct {
						Bold          bool   `json:"bold"`
						Italic        bool   `json:"italic"`
						Strikethrough bool   `json:"strikethrough"`
						Underline     bool   `json:"underline"`
						Code          bool   `json:"code"`
						Color         string `json:"color"`
					} `json:"annotations"`
					PlainText string `json:"plain_text"`
					Href      any    `json:"href"`
				} `json:"title"`
			} `json:"Name"`
			Due struct {
				ID   string `json:"id"`
				Type string `json:"type"`
				Date any    `json:"date"`
			} `json:"Due"`
		} `json:"properties"`
		URL string `json:"url"`
	} `json:"results"`
	NextCursor any    `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
	Type       string `json:"type"`
	Page       struct {
	} `json:"page"`
}

// parseNotionDate : Notionの日付プロパティ（日付のみ or 時刻付き）をパース
func parseNotionDate(s string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (dto NotionTaskDTO) ToTaskListDomainModel(ctx context.Context) task.List {
	tasks := make(task.List, 0, len(dto.Results))
	for _, r := range dto.Results {
		// タイトルが空のページはスキップ（panic 防止）
		if len(r.Properties.Name.Title) == 0 {
			continue
		}

		var deadline *time.Time
		if r.Properties.Deadline.Date.Start != "" {
			t, err := parseNotionDate(r.Properties.Deadline.Date.Start)
			if err != nil {
				log.Error(ctx, fmt.Errorf("Notion Deadlineのパースに失敗: id=%s start=%q: %w", r.ID, r.Properties.Deadline.Date.Start, err))
			} else {
				deadline = &t
			}
		}

		tasks = append(tasks, task.Task{
			ID:            r.ID,
			Title:         r.Properties.Name.Title[0].PlainText,
			Description:   "",
			Deadline:      deadline,
			Status:        task.Status(r.Properties.Status.Select.Name),
			List:          "All",
			ShortURL:      r.URL,
			OriginalModel: nil,
		})
	}
	return tasks
}

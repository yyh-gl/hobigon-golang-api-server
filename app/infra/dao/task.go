package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto/notion"
)

const defaultPageSize = 100

var notionHTTPClient = &http.Client{Timeout: 10 * time.Second}

type task struct {
	NotionAPIToken   string
	NotionDatabaseID string
}

// NewTask : タスク用のゲートウェイを取得
func NewTask() gateway.Task {
	return &task{
		NotionAPIToken:   os.Getenv("NOTION_API_TOKEN"),
		NotionDatabaseID: os.Getenv("NOTION_DATABASE_ID"),
	}
}

func (t task) fetchTasks(ctx context.Context, body notion.FetchTasksRequestBody) (model.List, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", t.NotionDatabaseID)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.NotionAPIToken))
	req.Header.Add("Notion-Version", "2022-06-28")

	res, err := notionHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var taskDTO notion.NotionTaskDTO
	if err := json.Unmarshal(resp, &taskDTO); err != nil {
		return nil, err
	}
	return taskDTO.ToTaskListDomainModel(), nil
}

// FetchActiveTasks : 『To Do』および『Doing』ステータスのタスクをすべて取得
func (t task) FetchActiveTasks(ctx context.Context) (model.List, error) {
	body := notion.FetchTasksRequestBody{
		PageSize: defaultPageSize,
		Filter: notion.OrFilter{
			Or: []notion.SingleFilter{
				{
					Property: "Status",
					Select:   &notion.Select{Equals: model.StatusToDo.String()},
				},
				{
					Property: "Status",
					Select:   &notion.Select{Equals: model.StatusDoing.String()},
				},
			},
		},
		Sorts: []notion.Sort{
			{Property: "Deadline", Direction: "ascending"},
			{Property: "Date Created", Direction: "ascending"},
		},
	}
	return t.fetchTasks(ctx, body)
}

// UpdateTaskStatus : 指定したタスクのステータスを更新
func (t task) UpdateTaskStatus(ctx context.Context, tsk model.Task, status model.Status) error {
	url := fmt.Sprintf("https://api.notion.com/v1/pages/%s", tsk.ID)

	body := map[string]any{
		"properties": map[string]any{
			"Status": map[string]any{
				"select": map[string]any{"name": status.String()},
			},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.NotionAPIToken))
	req.Header.Add("Notion-Version", "2022-06-28")

	res, err := notionHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()
	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("notion API returned status %d: %s", res.StatusCode, body)
	}
	return nil
}


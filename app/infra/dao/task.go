package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto/notion"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

const defaultPageSize = 100

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

// FetchCautionTasks : 今後1週間以内に期限が迫っているタスクを取得
// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨
func (t task) FetchCautionTasks(ctx context.Context) (model.List, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", t.NotionDatabaseID)

	body := notion.FetchTasksRequestBody{
		PageSize: defaultPageSize,
		Filter: notion.AndFilter{
			And: []notion.SingleFilter{
				{
					Property: "Deadline",
					Date: notion.Date{
						OnOrAfter: time.Now().Format("2006-01-02"),
					},
				},
				{
					Property: "Deadline",
					Date: notion.Date{
						OnOrBefore: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02"),
					},
				},
			},
		},
		Sorts: []notion.Sort{
			{Property: "Deadline", Direction: "ascending"},
			{Property: "Date Created", Direction: "ascending"},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.NotionAPIToken))
	req.Header.Add("Notion-Version", "2022-06-28")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
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

func (t task) FetchDeadTasks(ctx context.Context) (model.List, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/databases/%s/query", t.NotionDatabaseID)

	body := notion.FetchTasksRequestBody{
		PageSize: defaultPageSize,
		Filter: notion.SingleFilter{
			Property: "Deadline",
			Date: notion.Date{
				OnOrBefore: time.Now().Add(-1 * 24 * time.Hour).Format("2006-01-02"),
			},
		},
		Sorts: []notion.Sort{
			{Property: "Deadline", Direction: "ascending"},
			{Property: "Date Created", Direction: "ascending"},
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	payload := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.NotionAPIToken))
	req.Header.Add("Notion-Version", "2022-06-28")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
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

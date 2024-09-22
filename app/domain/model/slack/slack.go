package slack

import (
	"fmt"
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// Slack : Slackを表すドメインモデル
// TODO: ドメインモデル貧血症を治す
type Slack struct {
	Username string
	// TODO: enum化
	Channel string
	Text    string
}

// GetWebHookURL : WebHook URLを取得
func (s Slack) GetWebHookURL() (webHookURL string) {
	switch s.Channel {
	case "00_today_tasks":
		return os.Getenv("WEBHOOK_URL_TO_00")
	case "03_pokemon":
		return os.Getenv("WEBHOOK_URL_TO_03")
	case "51_tech_blog":
		return os.Getenv("WEBHOOK_URL_TO_51")
	case "2019新卒技術_雑談":
		return os.Getenv("WEBHOOK_URL_TO_SHINSOTSU")
	}
	return ""
}

// CreateTaskMessage : タスク通知用のメッセージを作成
// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨
func (s Slack) CreateTaskMessage(cautionTasks []task.Task, deadTasks []task.Task) string {
	message := ":mario2:Caution & ToDo Tasks:mario2:\n"
	for i, t := range cautionTasks {
		due := "なるはや"
		if t.Due != nil {
			due = t.Due.Format("2006-01-02")
		}
		message += fmt.Sprintf("%d: <%s|%s> :alarm_clock:`%s`\n", i+1, t.ShortURL, t.Title, due)
	}

	message += "\n\n:space_invader:Dead Tasks:space_invader:\n"
	for i, t := range deadTasks {
		message += fmt.Sprintf("%d: <%s|%s> :alarm_clock:`%s`\n", i+1, t.ShortURL, t.Title, t.Due.Format("2006-01-02"))
	}

	return message
}

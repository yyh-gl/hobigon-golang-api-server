package model

import (
	"os"
	"strconv"
)

type Slack struct {
	Username   string
	Channel    string
	Text       string
}

func (s Slack) GetWebHookURL() (webHookURL string) {
	switch s.Channel {
	case "00_today_tasks":
		webHookURL = os.Getenv("WEBHOOK_URL_TO_00")
	case "51_tech_blog":
		webHookURL = os.Getenv("WEBHOOK_URL_TO_51")
	}
	return webHookURL
}

func (s Slack) CreateTaskMessage(todayTasks []Task, dueOverTasks []Task) string {
	header := ":pencil2::pencil2::pencil2: 今日のタスク :pencil2::pencil2::pencil2:\n\n"
	header += "           :mega: 総タスク数 " + strconv.Itoa(len(todayTasks)) + " 個 :mega:\n"
	header += "   :space_invader: → 期限切れタスク数" + strconv.Itoa(len(dueOverTasks)) + "個 :space_invader:\n\n"

	body := "           :mario2::dash: :kata-me::kata-i::kata-nn: :mario2::dash:\n"
	// ここにタスクをメッセージ化していく処理を書く

	return header + body
}

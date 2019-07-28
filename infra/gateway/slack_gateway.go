package gateway

import (
	"fmt"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
)

type slackGateway struct {}

// TODO: 場所ここ？
func NewSlackGateway() gateway.SlackGateway {
	return &slackGateway{}
}

func (s slackGateway) send(data model.Slack) (err []error) {
	payload := slack.Payload{
		Username: data.Username,
		Channel: data.Channel,
		Text: data.Text,
	}

	webHookURL := data.GetWebHookURL()
	err = slack.Send(webHookURL, "", payload)
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return err
	}

	return nil
}

func (s slackGateway) SendTask(todayTasks []model.Task, dueOverTasks []model.Task) (err error) {
	data := model.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	s.send(data)
	return err
}

package gateway

import (
	"context"
	"log"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type slackGateway struct {}

// TODO: 場所ここ？
func NewSlackGateway() gateway.SlackGateway {
	return &slackGateway{}
}

func (s slackGateway) send(ctx context.Context, data model.Slack) (err []error) {
	logger := ctx.Value("logger").(log.Logger)

	payload := slack.Payload{
		Username: data.Username,
		Channel: data.Channel,
		Text: data.Text,
	}

	webHookURL := data.GetWebHookURL()
	err = slack.Send(webHookURL, "", payload)
	if err != nil {
		logger.Println(err)
		return err
	}

	return nil
}

func (s slackGateway) SendTask(ctx context.Context, todayTasks []model.Task, dueOverTasks []model.Task) (err error) {
	data := model.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	s.send(ctx, data)
	return err
}

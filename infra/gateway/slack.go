package gateway

import (
	"context"
	"strconv"

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
	payload := slack.Payload{
		Username: data.Username,
		Channel: data.Channel,
		Text: data.Text,
	}

	webHookURL := data.GetWebHookURL()
	err = slack.Send(webHookURL, "", payload)
	if err != nil {
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

func (s slackGateway) SendBirthday(ctx context.Context, birthday model.Birthday) (err error) {
	data := model.Slack{
		Username: "聖母マリア様",
		Channel:  "2019新卒技術_雑談",
	}

	data.Text = birthday.CreateBirthdayMessage()

	s.send(ctx, data)
	return err
}

func (s slackGateway) SendLikeNotify(ctx context.Context, blog model.Blog) (err error) {
	data := model.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text: "【" + blog.Title + "】いいね！（Total: " + strconv.Itoa(*blog.Count) + "）",
	}

	s.send(ctx, data)
	return err
}

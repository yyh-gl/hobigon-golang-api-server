package igateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/task"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/blog"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	mySlack "github.com/yyh-gl/hobigon-golang-api-server/domain/model/slack"
)

//////////////////////////////////////////////////
// NewSlackGateway
//////////////////////////////////////////////////

type slackGateway struct{}

// NewSlackGateway : Slack 用のゲートウェイを取得
func NewSlackGateway() gateway.SlackGateway {
	return &slackGateway{}
}

//////////////////////////////////////////////////
// send
//////////////////////////////////////////////////

// send : Slack に通知を送信
func (s slackGateway) send(ctx context.Context, data mySlack.Slack) (err []error) {
	payload := slack.Payload{
		Username: data.Username,
		Channel:  data.Channel,
		Text:     data.Text,
	}

	webHookURL := data.GetWebHookURL()
	err = slack.Send(webHookURL, "", payload)
	if err != nil {
		return err
	}

	return nil
}

//////////////////////////////////////////////////
// SendTask
//////////////////////////////////////////////////

// SendTask : Slack にタスクを送信
func (s slackGateway) SendTask(ctx context.Context, todayTasks []task.Task, dueOverTasks []task.Task) (err error) {
	data := mySlack.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	s.send(ctx, data)
	return err
}

//////////////////////////////////////////////////
// SendBirthday
//////////////////////////////////////////////////

// SendBirthday : Slack に誕生日通知を送信
func (s slackGateway) SendBirthday(ctx context.Context, birthday birthday.Birthday) (err error) {
	var data mySlack.Slack
	switch {
	case app.IsPrd():
		data = mySlack.Slack{
			Username: "聖母マリア様",
			Channel:  "2019新卒技術_雑談",
		}
	default:
		data = mySlack.Slack{
			Username: "まりお",
			Channel:  "00_today_tasks",
		}
	}

	data.Text = birthday.CreateBirthdayMessage()

	s.send(ctx, data)
	return err
}

//////////////////////////////////////////////////
// SendLikeNotify
//////////////////////////////////////////////////

// SendLikeNotify : Slack にいいね（ブログ）通知を送信
func (s slackGateway) SendLikeNotify(ctx context.Context, blog blog.Blog) (err error) {
	data := mySlack.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     blog.CreateLikeMessage(),
	}

	s.send(ctx, data)
	return err
}

//////////////////////////////////////////////////
// SendRanking
//////////////////////////////////////////////////

// SendRanking : Slack にアクセスラキングを送信
func (s slackGateway) SendRanking(ctx context.Context, ranking string) (err error) {
	data := mySlack.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     ranking,
	}

	s.send(ctx, data)
	return err
}

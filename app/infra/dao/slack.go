package dao

import (
	"context"

	slackWebHook "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	modelBd "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	modelB "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	modelS "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/slack"
	modelT "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

type slack struct{}

// NewSlack : Slack 用のゲートウェイを取得
func NewSlack() gateway.Slack {
	return &slack{}
}

// send : Slack に通知を送信
func (s slack) send(ctx context.Context, data modelS.Slack) (err []error) {
	payload := slackWebHook.Payload{
		Username: data.Username,
		Channel:  data.Channel,
		Text:     data.Text,
	}

	webHookURL := data.GetWebHookURL()
	err = slackWebHook.Send(webHookURL, "", payload)
	if err != nil {
		return err
	}

	return nil
}

// SendTask : Slack にタスクを送信
func (s slack) SendTask(ctx context.Context, todayTasks []modelT.Task, dueOverTasks []modelT.Task) (err error) {
	data := modelS.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	s.send(ctx, data)
	return err
}

// SendBirthday : Slack に誕生日通知を送信
func (s slack) SendBirthday(ctx context.Context, birthday modelBd.Birthday) (err error) {
	var data modelS.Slack
	switch {
	case app.IsPrd():
		data = modelS.Slack{
			Username: "聖母マリア様",
			Channel:  "2019新卒技術_雑談",
		}
	default:
		data = modelS.Slack{
			Username: "まりお",
			Channel:  "00_today_tasks",
		}
	}

	data.Text = birthday.CreateBirthdayMessage()

	s.send(ctx, data)
	return err
}

// SendLikeNotify : Slack にいいね（ブログ）通知を送信
func (s slack) SendLikeNotify(ctx context.Context, blog modelB.Blog) (err error) {
	data := modelS.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     blog.CreateLikeMessage(),
	}

	s.send(ctx, data)
	return err
}

// SendRanking : Slack にアクセスラキングを送信
func (s slack) SendRanking(ctx context.Context, ranking string) (err error) {
	data := modelS.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     ranking,
	}

	s.send(ctx, data)
	return err
}

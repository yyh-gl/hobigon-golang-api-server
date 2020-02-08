package dao

import (
	"context"
	"errors"

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
func (s slack) send(ctx context.Context, data modelS.Slack) error {
	payload := slackWebHook.Payload{
		Username: data.Username,
		Channel:  data.Channel,
		Text:     data.Text,
	}

	webHookURL := data.GetWebHookURL()
	errList := slackWebHook.Send(webHookURL, "", payload)
	if errList != nil {
		msg := ""
		for _, e := range errList {
			msg += e.Error() + " & "
		}
		return errors.New(msg)
	}
	return nil
}

// SendTask : Slack にタスクを送信
func (s slack) SendTask(ctx context.Context, todayTasks []modelT.Task, dueOverTasks []modelT.Task) error {
	data := modelS.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	if err := s.send(ctx, data); err != nil {
		return err
	}
	return nil
}

// SendBirthday : Slack に誕生日通知を送信
func (s slack) SendBirthday(ctx context.Context, birthdayList modelBd.BirthdayList) error {
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

	for _, birthday := range birthdayList {
		data.Text = birthday.CreateBirthdayMessage()
		if err := s.send(ctx, data); err != nil {
			return err
		}
	}
	return nil
}

// SendLikeNotify : Slack にいいね（ブログ）通知を送信
func (s slack) SendLikeNotify(ctx context.Context, blog modelB.Blog) error {
	data := modelS.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     blog.CreateLikeMessage(),
	}

	if err := s.send(ctx, data); err != nil {
		return err
	}
	return nil
}

// SendRanking : Slack にアクセスラキングを送信
func (s slack) SendRanking(ctx context.Context, ranking string) error {
	data := modelS.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     ranking,
	}

	if err := s.send(ctx, data); err != nil {
		return err
	}
	return nil
}

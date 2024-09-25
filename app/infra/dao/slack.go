package dao

import (
	"context"
	"errors"

	slackWebHook "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	modelB "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/pokemon"
	modelS "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/slack"
	modelT "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// TODO: It's strange that this file is under the dao package, so fix it.

type slack struct{}

// NewSlack : Slack用のゲートウェイを取得
func NewSlack() gateway.Slack {
	return &slack{}
}

// send : Slackに通知を送信
func (s slack) send(ctx context.Context, data modelS.Slack) error {
	payload := slackWebHook.Payload{
		Username: data.Username,
		Channel:  data.Channel,
		Text:     data.Text,
	}

	webHookURL := data.GetWebHookURL()
	errList := slackWebHook.Send(webHookURL, "", payload)
	if len(errList) != 0 && !app.IsTest() {
		msg := ""
		for _, e := range errList {
			if msg != "" {
				msg += " & "
			}
			msg += e.Error()
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

func (s slack) SendPokemonEvents(ctx context.Context, events []pokemon.Notification) error {
	text := ""
	if len(events) == 0 {
		text = "新しいイベント情報はありません。"
	} else {
		for _, e := range events {
			text += "・" + e.Title() + "\n"
		}
	}

	data := modelS.Slack{
		Username: "まりお",
		Channel:  "03_pokemon",
		Text:     text,
	}

	if err := s.send(ctx, data); err != nil {
		return err
	}
	return nil
}

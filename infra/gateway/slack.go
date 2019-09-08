package gateway

import (
	"strconv"

	"github.com/yyh-gl/hobigon-golang-api-server/app"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type slackGateway struct{}

func NewSlackGateway() gateway.SlackGateway {
	return &slackGateway{}
}

func (s slackGateway) send(data model.Slack) (err []error) {
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

func (s slackGateway) SendTask(todayTasks []model.Task, dueOverTasks []model.Task) (err error) {
	data := model.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	s.send(data)
	return err
}

func (s slackGateway) SendBirthday(birthday model.Birthday) (err error) {
	var data model.Slack
	switch {
	case app.IsPrd():
		data = model.Slack{
			Username: "聖母マリア様",
			Channel:  "2019新卒技術_雑談",
		}
	default:
		data = model.Slack{
			Username: "まりお",
			Channel:  "00_today_tasks",
		}
	}

	data.Text = birthday.CreateBirthdayMessage()

	s.send(data)
	return err
}

func (s slackGateway) SendLikeNotify(blog model.Blog) (err error) {
	data := model.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     "【" + blog.Title + "】いいね！（Total: " + strconv.Itoa(*blog.Count) + "）",
	}

	s.send(data)
	return err
}

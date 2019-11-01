package gateway

import (
	"strconv"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
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

//////////////////////////////////////////////////
// SendTask
//////////////////////////////////////////////////

// SendTask : Slack にタスクを送信
func (s slackGateway) SendTask(todayTasks []model.Task, dueOverTasks []model.Task) (err error) {
	data := model.Slack{
		Username: "まりお",
		Channel:  "00_today_tasks",
	}

	data.Text = data.CreateTaskMessage(todayTasks, dueOverTasks)

	s.send(data)
	return err
}

//////////////////////////////////////////////////
// SendBirthday
//////////////////////////////////////////////////

// SendBirthday : Slack に誕生日通知を送信
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

//////////////////////////////////////////////////
// SendLikeNotify
//////////////////////////////////////////////////

// SendLikeNotify : Slack にいいね（ブログ）通知を送信
func (s slackGateway) SendLikeNotify(blog model.Blog) (err error) {
	data := model.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     "【" + blog.Title + "】いいね！（Total: " + strconv.Itoa(*blog.Count) + "）",
	}

	s.send(data)
	return err
}

//////////////////////////////////////////////////
// SendRanking
//////////////////////////////////////////////////

// SendRanking : Slack にアクセスラキングを送信
func (s slackGateway) SendRanking(ranking string) (err error) {
	data := model.Slack{
		Username: "くりぼー",
		Channel:  "51_tech_blog",
		Text:     ranking,
	}

	s.send(data)
	return err
}

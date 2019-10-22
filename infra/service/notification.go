package service

import (
	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/service"
)

type notificationService struct {
	slackGateway gateway.SlackGateway
}

// NewNotificationService は NotificationService のインスタンスを生成
func NewNotificationService(slackGateway gateway.SlackGateway, taskGateway gateway.TaskGateway) service.NotificationService {
	return &notificationService{
		slackGateway: slackGateway,
	}
}

// SendBirthdayToSlack は今日の誕生日
func (n notificationService) SendTodayBirthdayToSlack(birthday model.Birthday) (err error) {
	// 今日が誕生日であった場合にのみ Slack に通知
	if birthday.IsToday() {
		if err = n.slackGateway.SendBirthday(birthday); err != nil {
			return errors.Wrap(err, "slackGateway.SendBirthday()内でのエラー")
		}
	}
	return nil
}

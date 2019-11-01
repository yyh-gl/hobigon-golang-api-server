package service

import (
	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/service"
)

//////////////////////////////////////////////////
// NewNotificationService
//////////////////////////////////////////////////

type notificationService struct {
	sg gateway.SlackGateway
}

// NewNotificationService : 通知用のサービスを取得
func NewNotificationService(sg gateway.SlackGateway) service.NotificationService {
	return &notificationService{
		sg: sg,
	}
}

// SendBirthdayToSlack は今日の誕生日
func (ns notificationService) SendTodayBirthdayToSlack(birthday model.Birthday) (err error) {
	// 今日が誕生日であった場合にのみ Slack に通知
	if birthday.IsToday() {
		if err = ns.sg.SendBirthday(birthday); err != nil {
			return errors.Wrap(err, "slackGateway.SendBirthday()内でのエラー")
		}
	}
	return nil
}

package iservice

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/service"
)

type notification struct {
	g gateway.SlackGateway
}

// NewNotification : Notification用ドメインサービスを取得
func NewNotification(g gateway.SlackGateway) service.Notification {
	return &notification{
		g: g,
	}
}

// SendBirthdayToSlackToSlack : 今日の誕生日を通知
func (n notification) SendTodayBirthdayToSlack(ctx context.Context, birthday model.Birthday) (err error) {
	// 今日が誕生日であった場合にのみ Slack に通知
	if birthday.Date().IsToday() {
		if err = n.g.SendBirthday(ctx, birthday); err != nil {
			return errors.Wrap(err, "slackGateway.SendBirthday()内でのエラー")
		}
	}
	return nil
}

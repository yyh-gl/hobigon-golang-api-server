package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/birthday"
)

// NotificationService : 通知用サービスのインターフェース
type NotificationService interface {
	SendTodayBirthdayToSlack(ctx context.Context, birthday birthday.Birthday) error
}

package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// NotificationService は通知用サービスのインターフェース
type NotificationService interface {
	SendTodayBirthdayToSlack(ctx context.Context, birthday model.Birthday) error
}

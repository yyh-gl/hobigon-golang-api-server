package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

// NotificationService : 通知用サービスのインターフェース
type NotificationService interface {
	SendTodayBirthdayToSlack(ctx context.Context, birthday entity.Birthday) error
}

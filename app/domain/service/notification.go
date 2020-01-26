package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

// Notification : Notification用ドメインサービスのインターフェース
type Notification interface {
	SendTodayBirthdayToSlack(ctx context.Context, birthday birthday.Birthday) error
}

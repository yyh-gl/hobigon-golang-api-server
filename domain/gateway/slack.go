package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// SlackGateway : 通知用のゲートウェイインターフェース
type SlackGateway interface {
	SendTask(ctx context.Context, todayTasks []model.Task, dueOverTasks []model.Task) error
	SendBirthday(ctx context.Context, birthday model.Birthday) error
	SendLikeNotify(ctx context.Context, blog model.Blog) error
	SendRanking(ctx context.Context, rankin string) error
}

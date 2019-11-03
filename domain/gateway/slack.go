package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

// SlackGateway : 通知用のゲートウェイインターフェース
type SlackGateway interface {
	SendTask(ctx context.Context, todayTasks []entity.Task, dueOverTasks []entity.Task) error
	SendBirthday(ctx context.Context, birthday entity.Birthday) error
	SendLikeNotify(ctx context.Context, blog entity.Blog) error
	SendRanking(ctx context.Context, rankin string) error
}

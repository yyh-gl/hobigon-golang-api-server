package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// SlackGateway : 通知用のゲートウェイインターフェース
type SlackGateway interface {
	SendTask(ctx context.Context, todayTasks []task.Task, dueOverTasks []task.Task) error
	SendBirthday(ctx context.Context, birthday birthday.Birthday) error
	SendLikeNotify(ctx context.Context, blog blog.Blog) error
	SendRanking(ctx context.Context, rankin string) error
}

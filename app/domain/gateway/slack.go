package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/blog"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// Slack : Slack用ゲートウェイのインターフェース
type Slack interface {
	SendTask(ctx context.Context, todayTasks []task.Task, dueOverTasks []task.Task) error
	SendLikeNotify(ctx context.Context, blog blog.Blog) error
	SendRanking(ctx context.Context, rankin string) error
}

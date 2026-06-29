package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// Task : Task用ゲートウェイのインターフェース
type Task interface {
	FetchActiveTasks(context.Context) (task.List, error)
	UpdateTaskStatus(context.Context, task.Task, task.Status) error
}

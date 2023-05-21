package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// Task : Task用ゲートウェイのインターフェース
// FIXME: Trello -> Notion への移行を突貫工事で作ったのでリファクタ推奨
type Task interface {
	FetchCautionAndToDoTasks(context.Context) (task.List, error)
	FetchDeadTasks(context.Context) (task.List, error)
}

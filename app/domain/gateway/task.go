package gateway

import (
	"context"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

// Task : Task用ゲートウェイのインターフェース
type Task interface {
	// ここだけ処理簡易化のために外部ライブラリ（trello）への依存を許す
	GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error)
	GetTasksFromList(ctx context.Context, list trello.List) (task.List, task.List, error)
	MoveToWIP(ctx context.Context, tasks []task.Task) (err error)
}

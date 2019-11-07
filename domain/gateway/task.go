package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/task"

	"github.com/adlio/trello"
)

type TaskGateway interface {
	GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error)
	GetTasksFromList(ctx context.Context, list trello.List) (task.TaskList, task.TaskList, error)
	MoveToWIP(ctx context.Context, tasks []task.Task) (err error)
}

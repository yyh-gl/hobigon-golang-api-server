package gateway

import (
	"context"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type TaskGateway interface {
	GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error)
	GetTasksFromList(ctx context.Context, list trello.List) (model.TaskList, model.TaskList, error)
	MoveToWIP(ctx context.Context, tasks []model.Task) (err error)
}

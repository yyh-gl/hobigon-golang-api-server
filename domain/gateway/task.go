package gateway

import (
	"context"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

type TaskGateway interface {
	GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error)
	GetTasksFromList(ctx context.Context, list trello.List) (entity.TaskList, entity.TaskList, error)
	MoveToWIP(ctx context.Context, tasks []entity.Task) (err error)
}

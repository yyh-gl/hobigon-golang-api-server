package repository

import (
	"context"
	"github.com/adlio/trello"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type TaskRepository interface {
	GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error)
	GetTasksFromList(ctx context.Context, list trello.List) ([]*model.Task, error)
}

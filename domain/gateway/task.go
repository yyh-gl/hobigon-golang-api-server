package gateway

import (
	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type TaskGateway interface {
	GetListsByBoardID(boardID string) (lists []*trello.List, err error)
	GetTasksFromList(list trello.List) (model.TaskList, model.TaskList, error)
}

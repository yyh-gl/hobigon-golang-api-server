package repository

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type TaskRepository interface {
	GetBoardByID(context.Context) (*model.Board, error)
}

package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type SlackGateway interface {
	SendTask(ctx context.Context, todayTasks []model.Task, dueOverTasks []model.Task) error
}

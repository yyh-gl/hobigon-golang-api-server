package gateway

import "github.com/yyh-gl/hobigon-golang-api-server/domain/model"

type SlackGateway interface {
	SendTask(todayTasks []model.Task, dueOverTasks []model.Task) error
}

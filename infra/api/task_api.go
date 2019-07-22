package api

import (
	"context"

	//"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type taskRepository struct {}

func NewTaskRepository() repository.TaskRepository {
	return &taskRepository{}
}

func (ta taskRepository) GetTodayTasks(ctx context.Context) []*model.Task {
	var taskList []*model.Task
	return taskList
}

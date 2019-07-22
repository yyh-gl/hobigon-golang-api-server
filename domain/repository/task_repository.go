package repository

import (
"net/http"

"github.com/julienschmidt/httprouter"
"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type TaskRepository interface {
	GetTodayTasks(http.ResponseWriter, *http.Request, httprouter.Params) []*model.Task
}

type taskRepository struct {}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

func (ta taskRepository) GetTodayTasks(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) []*model.Task {
	var taskList []*model.Task
	return taskList
}

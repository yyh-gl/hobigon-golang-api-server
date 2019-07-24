package api

import (
	"context"
	"fmt"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type taskRepository struct {}

func NewTaskRepository() repository.TaskRepository {
	return &taskRepository{}
}

func (ta taskRepository) GetTodayTasks(ctx context.Context) (taskList []*model.Task, err error) {
	client := trello.NewClient("4d8b29f534d118c7abe797f84e819213", "044b77489fee00d98c8ab7ee91c0d23d4deb0ec99083794a2363587f381c2edc")
	board, err := client.GetBoard("qU5LuAv3", trello.Defaults())
	if err != nil {
		return taskList, err
	}

	fmt.Println("========================")
	fmt.Println(board)
	fmt.Println("========================")

	return taskList, nil
}

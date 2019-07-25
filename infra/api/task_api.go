package api

import (
	"context"
	"fmt"
	"os"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/repository"
)

type taskRepository struct {
	ApiKey   string
	ApiToken string
	MainBoardID string
}

func NewTaskRepository() repository.TaskRepository {
	return &taskRepository{
		ApiKey:   os.Getenv("TRELLO_API_KEY"),
		ApiToken: os.Getenv("TRELLO_API_TOKEN"),
	}
}

func (tr taskRepository) getBoard(ctx context.Context, boardID string) (board *trello.Board, err error) {
	client := trello.NewClient(tr.ApiKey, tr.ApiToken)
	board, err = client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return nil, err
	}
	return board, nil
}

func (tr taskRepository) GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error) {
	board, err := tr.getBoard(ctx, boardID)
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return nil, err
	}

	lists, err = board.GetLists(trello.Defaults())
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return nil, err
	}

	return lists, nil
}

func (tr taskRepository) GetTasksFromList(ctx context.Context, list trello.List) (tasks []*model.Task, err error) {
	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return nil, err
	}

	tasks = model.ConvertToTasksModel(trelloTasks)
	return tasks, nil
}

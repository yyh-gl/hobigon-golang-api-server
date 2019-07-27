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

// 場所ここ？ + gateway にかえる
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

	tasks = convertToTasksModel(trelloTasks)
	return tasks, nil
}

func convertToTasksModel(trelloCards []*trello.Card) (tasks []*model.Task) {
	for _, card := range trelloCards {
		task := new(model.Task)
		task.Title       = card.Name
		task.Description = card.Desc
		if card.Due != nil {
			task.Due = task.GetJSTDue(*card.Due)
		}
		tasks = append(tasks, task)
	}
	return tasks
}

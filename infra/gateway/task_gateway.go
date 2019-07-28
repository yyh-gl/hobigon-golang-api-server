package gateway

import (
	"context"
	"fmt"
	"os"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type taskGateway struct {
	ApiKey   string
	ApiToken string
	MainBoardID string
}

// TODO: 場所ここ？
func NewTaskGateway() gateway.TaskGateway {
	return &taskGateway{
		ApiKey:   os.Getenv("TRELLO_API_KEY"),
		ApiToken: os.Getenv("TRELLO_API_TOKEN"),
	}
}

func (tr taskGateway) getBoard(ctx context.Context, boardID string) (board *trello.Board, err error) {
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

func (tr taskGateway) GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error) {
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

func (tr taskGateway) GetTasksFromList(ctx context.Context, list trello.List) (taskList model.TaskList, dueOverTaskList model.TaskList, err error) {
	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return model.TaskList{}, model.TaskList{}, err
	}

	allTask := convertToTasksModel(trelloTasks)

	for _, task := range allTask.Tasks {
		if task.Due != nil && task.IsDueOver() {
			dueOverTaskList.Tasks = append(dueOverTaskList.Tasks, task)
		} else {
			taskList.Tasks = append(taskList.Tasks, task)
		}
	}
	return taskList, dueOverTaskList, nil
}

func convertToTasksModel(trelloCards []*trello.Card) (taskList model.TaskList) {
	for _, card := range trelloCards {
		task := new(model.Task)
		task.Title       = card.Name
		task.Description = card.Desc
		if card.Due != nil {
			task.Due = task.GetJSTDue(card.Due)
		}
		taskList.Tasks = append(taskList.Tasks, *task)
	}
	return taskList
}

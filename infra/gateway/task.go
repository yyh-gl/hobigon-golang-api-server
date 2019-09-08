package gateway

import (
	"os"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type taskGateway struct {
	APIKey   string
	APIToken string
}

func NewTaskGateway() gateway.TaskGateway {
	return &taskGateway{
		APIKey:   os.Getenv("TRELLO_API_KEY"),
		APIToken: os.Getenv("TRELLO_API_TOKEN"),
	}
}

func (tr taskGateway) getBoard(boardID string) (board *trello.Board, err error) {
	client := trello.NewClient(tr.APIKey, tr.APIToken)
	board, err = client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (tr taskGateway) GetListsByBoardID(boardID string) (lists []*trello.List, err error) {
	board, err := tr.getBoard(boardID)
	if err != nil {
		return nil, err
	}

	// TODO: ここで todo と wip だけにしちゃう
	lists, err = board.GetLists(trello.Defaults())
	if err != nil {
		return nil, err
	}

	// Board情報付与
	for _, list := range lists {
		list.Board = board
	}

	return lists, nil
}

// TODO: WIP 外にあるやつで期限が今日のものを WIP に移動させる機能 追加
func (tr taskGateway) GetTasksFromList(list trello.List) (taskList model.TaskList, dueOverTaskList model.TaskList, err error) {
	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		return model.TaskList{}, model.TaskList{}, err
	}

	allTask := convertToTasksModel(trelloTasks)

	for _, task := range allTask.Tasks {
		task.Board = list.Board.Name
		task.List = list.Name

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
		task.Title = card.Name
		task.Description = card.Desc
		task.ShortURL = card.ShortURL
		if card.Due != nil {
			task.Due = task.GetJSTDue(card.Due)
		}
		taskList.Tasks = append(taskList.Tasks, *task)
	}
	return taskList
}

package gateway

import (
	"context"
	"log"
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
	logger := ctx.Value("logger").(log.Logger)

	client := trello.NewClient(tr.ApiKey, tr.ApiToken)
	board, err = client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return board, nil
}

func (tr taskGateway) GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error) {
	logger := ctx.Value("logger").(log.Logger)

	board, err := tr.getBoard(ctx, boardID)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// TODO: ここで todo と wip だけにしちゃう
	lists, err = board.GetLists(trello.Defaults())
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	// Board情報付与
	for _, list := range lists {
		list.Board = board
	}

	return lists, nil
}

func (tr taskGateway) GetTasksFromList(ctx context.Context, list trello.List) (taskList model.TaskList, dueOverTaskList model.TaskList, err error) {
	logger := ctx.Value("logger").(log.Logger)

	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		logger.Println(err)
		return model.TaskList{}, model.TaskList{}, err
	}

	allTask := convertToTasksModel(trelloTasks)

	for _, task := range allTask.Tasks {
		task.Board = list.Board.Name
		task.List  = list.Name

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
		task.ShortURL    = card.ShortURL
		if card.Due != nil {
			task.Due = task.GetJSTDue(card.Due)
		}
		taskList.Tasks = append(taskList.Tasks, *task)
	}
	return taskList
}

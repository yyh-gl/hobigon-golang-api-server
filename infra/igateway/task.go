package igateway

import (
	"context"
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/task"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
)

//////////////////////////////////////////////////
// NewTaskGateway
//////////////////////////////////////////////////

type taskGateway struct {
	APIKey   string
	APIToken string
}

// NewTaskGateway : タスク用のゲートウェイを取得
func NewTaskGateway() gateway.TaskGateway {
	return &taskGateway{
		APIKey:   os.Getenv("TRELLO_API_KEY"),
		APIToken: os.Getenv("TRELLO_API_TOKEN"),
	}
}

//////////////////////////////////////////////////
// getBoard
//////////////////////////////////////////////////

// getBoard : ボード情報を取得
func (tg taskGateway) getBoard(ctx context.Context, boardID string) (board *trello.Board, err error) {
	client := trello.NewClient(tg.APIKey, tg.APIToken)
	board, err = client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return board, nil
}

//////////////////////////////////////////////////
// GetListsByBoardID
//////////////////////////////////////////////////

// GetListsByBoardID : ボードIDからリスト情報を取得
func (tg taskGateway) GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error) {
	board, err := tg.getBoard(ctx, boardID)
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

//////////////////////////////////////////////////
// GetTasksFromList
//////////////////////////////////////////////////

// GetTasksFromList : リストからタスク一覧を取得
func (tg taskGateway) GetTasksFromList(ctx context.Context, list trello.List) (taskList task.TaskList, dueOverTaskList task.TaskList, err error) {
	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		return task.TaskList{}, task.TaskList{}, err
	}

	allTask := convertToTasksModel(ctx, trelloTasks)

	for _, t := range allTask.Tasks {
		t.Board = list.Board.Name
		t.List = list.Name

		// 期限切れタスクを抽出
		if t.Due != nil && t.IsDueOver() {
			dueOverTaskList.Tasks = append(dueOverTaskList.Tasks, t)
		} else {
			taskList.Tasks = append(taskList.Tasks, t)
		}
	}
	return taskList, dueOverTaskList, nil
}

//////////////////////////////////////////////////
// convertToTasksModel
//////////////////////////////////////////////////

// convertToTasksModel : infra 層用の Task モデルをドメインモデルに変換
func convertToTasksModel(ctx context.Context, trelloCards []*trello.Card) (taskList task.TaskList) {
	for _, card := range trelloCards {
		t := new(task.Task)
		t.Title = card.Name
		t.Description = card.Desc
		t.ShortURL = card.ShortURL
		if card.Due != nil {
			t.Due = t.GetJSTDue(card.Due)
		}
		t.OriginalModel = card
		taskList.Tasks = append(taskList.Tasks, *t)
	}
	return taskList
}

//////////////////////////////////////////////////
// MoveToWIP
//////////////////////////////////////////////////

// MoveToWIP : 指定タスクを WIP リストに移動
func (tg taskGateway) MoveToWIP(ctx context.Context, tasks []task.Task) (err error) {
	for _, t := range tasks {
		var wipListID string
		switch t.Board {
		case "Main":
			wipListID = os.Getenv("MAIN_WIP_LIST_ID")
		case "Tech":
			wipListID = os.Getenv("TECH_WIP_LIST_ID")
		case "Work":
			wipListID = os.Getenv("WORK_WIP_LIST_ID")
		}

		card := t.OriginalModel.(*trello.Card)
		err = card.MoveToList(wipListID, trello.Defaults())
		if err != nil {
			// TODO: DB操作ではないので、途中で失敗した場合ロールバックできない問題を考える
			return err
		}
	}
	return nil
}

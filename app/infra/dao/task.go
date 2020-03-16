package dao

import (
	"context"
	"os"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/task"
)

type task struct {
	APIKey   string
	APIToken string
}

// NewTask : タスク用のゲートウェイを取得
func NewTask() gateway.Task {
	return &task{
		APIKey:   os.Getenv("TRELLO_API_KEY"),
		APIToken: os.Getenv("TRELLO_API_TOKEN"),
	}
}

// getBoard : ボード情報を取得
func (t task) getBoard(ctx context.Context, boardID string) (board *trello.Board, err error) {
	client := trello.NewClient(t.APIKey, t.APIToken)
	board, err = client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return board, nil
}

// GetListsByBoardID : ボードIDからリスト情報を取得
func (t task) GetListsByBoardID(ctx context.Context, boardID string) (lists []*trello.List, err error) {
	board, err := t.getBoard(ctx, boardID)
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

// GetTasksFromList : リストからタスク一覧を取得
func (t task) GetTasksFromList(ctx context.Context, list trello.List) (taskList model.List, dueOverTaskList model.List, err error) {
	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		return model.List{}, model.List{}, err
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

// convertToTasksModel : infra層用のTaskモデルをドメインモデルに変換
func convertToTasksModel(ctx context.Context, trelloCards []*trello.Card) (taskList model.List) {
	for _, card := range trelloCards {
		t := model.Task{}
		t.Title = card.Name
		t.Description = card.Desc
		t.ShortURL = card.ShortURL
		if card.Due != nil {
			t.Due = t.GetJSTDue(card.Due)
		}
		t.OriginalModel = card
		taskList.Tasks = append(taskList.Tasks, t)
	}
	return taskList
}

// MoveToWIP : 指定タスクをWIPリストに移動
func (t task) MoveToWIP(ctx context.Context, tasks []model.Task) (err error) {
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

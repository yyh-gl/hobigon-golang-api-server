package gateway

import (
	"os"

	"github.com/adlio/trello"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/gateway"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
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
func (tg taskGateway) getBoard(boardID string) (board *trello.Board, err error) {
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
func (tg taskGateway) GetListsByBoardID(boardID string) (lists []*trello.List, err error) {
	board, err := tg.getBoard(boardID)
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
func (tg taskGateway) GetTasksFromList(list trello.List) (taskList model.TaskList, dueOverTaskList model.TaskList, err error) {
	trelloTasks, err := list.GetCards(trello.Defaults())
	if err != nil {
		return model.TaskList{}, model.TaskList{}, err
	}

	allTask := convertToTasksModel(trelloTasks)

	for _, task := range allTask.Tasks {
		task.Board = list.Board.Name
		task.List = list.Name

		// 期限切れタスクを抽出
		if task.Due != nil && task.IsDueOver() {
			dueOverTaskList.Tasks = append(dueOverTaskList.Tasks, task)
		} else {
			taskList.Tasks = append(taskList.Tasks, task)
		}
	}
	return taskList, dueOverTaskList, nil
}

//////////////////////////////////////////////////
// convertToTasksModel
//////////////////////////////////////////////////

// convertToTasksModel : infra 層用の Task モデルをドメインモデルに変換
func convertToTasksModel(trelloCards []*trello.Card) (taskList model.TaskList) {
	for _, card := range trelloCards {
		task := new(model.Task)
		task.Title = card.Name
		task.Description = card.Desc
		task.ShortURL = card.ShortURL
		if card.Due != nil {
			task.Due = task.GetJSTDue(card.Due)
		}
		task.OriginalModel = card
		taskList.Tasks = append(taskList.Tasks, *task)
	}
	return taskList
}

//////////////////////////////////////////////////
// MoveToWIP
//////////////////////////////////////////////////

// MoveToWIP : 指定タスクを WIP リストに移動
func (tg taskGateway) MoveToWIP(tasks []model.Task) (err error) {
	for _, task := range tasks {
		var wipListID string
		switch task.Board {
		case "Main":
			wipListID = os.Getenv("MAIN_WIP_LIST_ID")
		case "Tech":
			wipListID = os.Getenv("TECH_WIP_LIST_ID")
		case "Work":
			wipListID = os.Getenv("WORK_WIP_LIST_ID")
		}

		card := task.OriginalModel.(*trello.Card)
		err = card.MoveToList(wipListID, trello.Defaults())
		if err != nil {
			// TODO: DB操作ではないので、途中で失敗した場合ロールバックできない問題を考える
			return err
		}
	}
	return nil
}

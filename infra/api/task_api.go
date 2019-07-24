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
		MainBoardID: os.Getenv("MAIN_BOARD_ID"),
	}
}

func (tr taskRepository) GetBoardByID(ctx context.Context) (board *model.Board, err error) {
	client := trello.NewClient(tr.ApiKey, tr.ApiToken)
	trelloBoard, err := client.GetBoard(tr.MainBoardID, trello.Defaults())
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return nil, err
	}

	board = model.ConvertToBoardModel(*trelloBoard)

	return board, nil
}

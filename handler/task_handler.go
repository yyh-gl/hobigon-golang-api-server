package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/api"
)

type response struct {
	Success  string
	TaskList []*model.Task
}

func TaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: コンテキストの設定方法・場所のベストプラクティスが分かり次第修正
	ctx := r.Context()
	ctx = context.WithValue(ctx, "params", ps)

	taskApi := api.NewTaskRepository()
	lists, err := taskApi.GetListsByBoardID(ctx, os.Getenv("MAIN_BOARD_ID"))
	if err != nil {
		// TODO: ロガーに差し替え
		fmt.Println("v===== ERROR =====v")
		fmt.Println(err)
		fmt.Println("^===== ERROR =====^")
		return
	}

	var tasks []*model.Task
	for _, list := range lists {
		if list.Name == "ToDo(Private)" {
			tasks, err = taskApi.GetTasksFromList(ctx, *list)
			if err != nil {
				// TODO: ロガーに差し替え
				fmt.Println("v===== ERROR =====v")
				fmt.Println(err)
				fmt.Println("^===== ERROR =====^")
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response{ Success: "ok", TaskList: tasks}); err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

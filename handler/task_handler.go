package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/api"
)

type response struct {
	Success  string
	TaskList []*model.Task
}

func TaskHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	ctx := context.Background()

	taskApi := api.NewTaskRepository()
	taskList := taskApi.GetTodayTasks(ctx)
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response{ Success: "ok", TaskList: taskList}); err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

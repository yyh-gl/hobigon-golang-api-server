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

func TaskHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: コンテキストの設定方法・場所のベストプラクティスが分かり次第修正
	ctx := r.Context()
	ctx = context.WithValue(ctx, "params", ps)

	taskApi := api.NewTaskRepository()
	taskList, _ := taskApi.GetTodayTasks(ctx)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response{ Success: "ok", TaskList: taskList}); err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

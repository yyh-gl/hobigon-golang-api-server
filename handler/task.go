package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/app"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
)

func NotifyTaskHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		TaskList        []model.Task `json:"task_list"`
		DueOverTaskList []model.Task `json:"due_over_task_list"`
	}

	logger := app.Logger

	taskGateway := gateway.NewTaskGateway()
	slackGateway := gateway.NewSlackGateway()

	var todayTasks []model.Task
	var dueOverTasks []model.Task
	// TODO: ハンドラーにロジックもっちゃっているのを直したみ
	boardIDList := [3]string{os.Getenv("MAIN_BOARD_ID"), os.Getenv("TECH_BOARD_ID"), os.Getenv("WORK_BOARD_ID")}
	for _, boardID := range boardIDList {
		lists, err := taskGateway.GetListsByBoardID(boardID)
		if err != nil {
			logger.Println(err)
			return
		}

		for _, list := range lists {
			// TODO: 今後必要があれば動的に変更できる仕組みを追加
			if list.Name == "TODO" || list.Name == "WIP" {
				taskList, dueOverTaskList, err := taskGateway.GetTasksFromList(*list)
				if err != nil {
					logger.Println(err)
					return
				}

				switch list.Name {
				case "TODO":
					// TODOリストからは今日のタスクのみ出力
					tasks := taskList.GetTodayTasks()
					todayTasks = append(todayTasks, tasks...)
				case "WIP":
					// WIPリストにあるタスクは全て出力
					todayTasks = append(todayTasks, taskList.Tasks...)
				}

				// 期限切れタスクは問答無用で通知
				dueOverTasks = append(dueOverTasks, dueOverTaskList.Tasks...)
			}
		}
	}

	err := slackGateway.SendTask(todayTasks, dueOverTasks)
	if err != nil {
		logger.Println(err)
		return
	}

	res := response{
		TaskList:        todayTasks,
		DueOverTaskList: dueOverTasks,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

type notificationResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// NotifyTodayTaskToSlackHandler は今日のタスク一覧を Slack に通知
func NotifyTodayTasksToSlackHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := notificationResponse{
		OK: true,
	}

	if err := usecase.NotifyTodayTasksToSlackUseCase(r.Context()); err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

// NotifyTodayBirthdayToSlackHandler は今日誕生日の人を Slack に通知
func NotifyTodayBirthdayToSlackHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := notificationResponse{
		OK: true,
	}

	if err := usecase.NotifyTodayBirthdayToSlackUseCase(r.Context()); err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

// NotifyAccessRankingHandler はアクセスランキングを Slack に通知
// TODO: robo タスクとしても実行できるようにしたい
func NotifyAccessRankingToSlackHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := notificationResponse{
		OK: true,
	}

	if err := usecase.NotifyAccessRankingUseCase(r.Context()); err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

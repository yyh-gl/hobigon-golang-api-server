package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/usecase"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

type response struct {
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error,omitempty"`
}

// NotifyTodayBirthdayToSlackHandler は今日誕生日の人を Slack に通知
func NotifyTodayBirthdayToSlackHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := response{
		IsSuccess: true,
	}

	if err := usecase.NotifyTodayBirthdayToSlackUseCase(r.Context()); err != nil {
		logger.Println(err)

		res.IsSuccess = false
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
func NotifyAccessRankingHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := response{
		IsSuccess: true,
	}

	if err := usecase.NotifyAccessRankingUseCase(r.Context()); err != nil {
		logger.Println(err)

		res.IsSuccess = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

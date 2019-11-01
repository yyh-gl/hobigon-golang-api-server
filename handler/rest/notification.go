package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

//////////////////////////////////////////////////
// NewNotificationHandler
//////////////////////////////////////////////////

// NotificationHandler : Slack 通知用のハンドラーインターフェース
type SlackNotificationHandler interface {
	NotifyTodayTasks(w http.ResponseWriter, r *http.Request)
	NotifyTodayBirthday(w http.ResponseWriter, r *http.Request)
	NotifyAccessRanking(w http.ResponseWriter, r *http.Request)
}

type slackNotificationHandler struct {
	nu usecase.NotificationUseCase
}

// NewNotificationHandler : Slack 通知用のハンドラーを取得
func NewSlackNotificationHandler(nu usecase.NotificationUseCase) SlackNotificationHandler {
	return &slackNotificationHandler{
		nu: nu,
	}
}

//////////////////////////////////////////////////
// 通知系共通処理
//////////////////////////////////////////////////

// 通知系 API の共通レスポンス
type notificationResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

//////////////////////////////////////////////////
// NotifyTodayTask
//////////////////////////////////////////////////

// NotifyTodayTaskToSlack : 今日のタスク一覧を Slack に通知
func (snh slackNotificationHandler) NotifyTodayTasks(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := notificationResponse{
		OK: true,
	}

	if err := snh.nu.NotifyTodayTasksToSlack(r.Context()); err != nil {
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

//////////////////////////////////////////////////
// NotifyTodayBirthday
//////////////////////////////////////////////////

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (snh slackNotificationHandler) NotifyTodayBirthday(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := notificationResponse{
		OK: true,
	}

	if err := snh.nu.NotifyTodayBirthdayToSlack(r.Context()); err != nil {
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

//////////////////////////////////////////////////
// NotifyAccessRanking
//////////////////////////////////////////////////

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (snh slackNotificationHandler) NotifyAccessRanking(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	res := notificationResponse{
		OK: true,
	}

	if err := snh.nu.NotifyAccessRanking(r.Context()); err != nil {
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

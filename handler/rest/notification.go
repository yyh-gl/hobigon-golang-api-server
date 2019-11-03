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
type NotificationHandler interface {
	NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request)
	NotifyTodayBirthdayToSlack(w http.ResponseWriter, r *http.Request)
	NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request)
}

type notificationHandler struct {
	nu usecase.NotificationUseCase
}

// NewNotificationHandler : Slack 通知用のハンドラーを取得
func NewNotificationHandler(nu usecase.NotificationUseCase) NotificationHandler {
	return &notificationHandler{
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
// NotifyTodayTasksToSlack
//////////////////////////////////////////////////

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (snh notificationHandler) NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request) {
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
// NotifyTodayBirthdayToSlack
//////////////////////////////////////////////////

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (snh notificationHandler) NotifyTodayBirthdayToSlack(w http.ResponseWriter, r *http.Request) {
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
// NotifyAccessRankingToSlack
//////////////////////////////////////////////////

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (snh notificationHandler) NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request) {
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

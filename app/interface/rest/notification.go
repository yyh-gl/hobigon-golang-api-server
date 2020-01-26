package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Notification : Notification用REST Handlerのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request)
	NotifyTodayBirthdayToSlack(w http.ResponseWriter, r *http.Request)
	NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request)
}

type notification struct {
	u usecase.Notification
}

// NewNotification : Notification用REST Handlerを取得
func NewNotification(u usecase.Notification) Notification {
	return &notification{
		u: u,
	}
}

// notificationResponse : Birthday用共通レスポンス
// TODO: OK, Error 部分は共通レスポンスにする
type notificationResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (n notification) NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request) {
	res := notificationResponse{
		OK: true,
	}

	if err := n.u.NotifyTodayTasksToSlack(r.Context()); err != nil {
		app.Logger.Println(err)

		res.OK = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (n notification) NotifyTodayBirthdayToSlack(w http.ResponseWriter, r *http.Request) {
	res := notificationResponse{
		OK: true,
	}

	if err := n.u.NotifyTodayBirthdayToSlack(r.Context()); err != nil {
		app.Logger.Println(err)

		res.OK = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (n notification) NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request) {
	res := notificationResponse{
		OK: true,
	}

	if err := n.u.NotifyAccessRanking(r.Context()); err != nil {
		app.Logger.Println(err)

		res.OK = false
		res.Error = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		app.Logger.Println(err)
		http.Error(w, "API レスポンスの JSON エンコードに失敗", http.StatusInternalServerError)
		return
	}
}

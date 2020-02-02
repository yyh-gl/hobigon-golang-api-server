package rest

import (
	"fmt"
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
	errorResponse
}

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (n notification) NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request) {
	resp := notificationResponse{}
	if err := n.u.NotifyTodayTasksToSlack(r.Context()); err != nil {
		errInfo := fmt.Errorf("notificationUseCase.NotifyTodayTasksToSlack()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}

	DoResponse(w, resp, http.StatusOK)
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (n notification) NotifyTodayBirthdayToSlack(w http.ResponseWriter, r *http.Request) {
	resp := notificationResponse{}
	if err := n.u.NotifyTodayBirthdayToSlack(r.Context()); err != nil {
		errInfo := fmt.Errorf("notificationUseCase.NotifyTodayBirthdayToSlack()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
	}

	DoResponse(w, resp, http.StatusOK)
}

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (n notification) NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request) {
	resp := notificationResponse{}
	if err := n.u.NotifyAccessRanking(r.Context()); err != nil {
		errInfo := fmt.Errorf("notificationUseCase.NotifyAccessRanking()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
	}

	DoResponse(w, resp, http.StatusOK)
}

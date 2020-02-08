package rest

import (
	"errors"
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

// notificationResponse : Notification用共通レスポンス
type notificationResponse struct {
	NotifiedNum int `json:"notified_num,omitempty"`
	errorResponse
}

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (n notification) NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request) {
	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyTodayTasksToSlack(r.Context())
	if err != nil {
		errInfo := fmt.Errorf("notificationUseCase.NotifyTodayTasksToSlack()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(w, resp, http.StatusOK)
}

// NotifyTodayBirthdayToSlack : 今日誕生日の人を Slack に通知
func (n notification) NotifyTodayBirthdayToSlack(w http.ResponseWriter, r *http.Request) {
	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyTodayBirthdayToSlack(r.Context())
	if err != nil {
		errInfo := fmt.Errorf("notificationUseCase.NotifyTodayBirthdayToSlack()でエラー: %w", err)
		app.Logger.Println(errInfo)

		if errors.Is(err, usecase.ErrBirthdayNotFound) {
			DoResponse(w, resp, http.StatusOK)
			return
		}

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(w, resp, http.StatusOK)
}

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (n notification) NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request) {
	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyAccessRanking(r.Context())
	if err != nil {
		errInfo := fmt.Errorf("notificationUseCase.NotifyAccessRanking()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(w, resp, http.StatusOK)
}

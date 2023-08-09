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
	NotifiedNum int `json:"notified_num"`
}

// NotifyTodayTasksToSlack : 今日のタスク一覧を Slack に通知
func (n notification) NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyTodayTasksToSlack(r.Context())
	if err != nil {
		app.Error(ctx, fmt.Errorf("notificationUseCase.NotifyTodayTasksToSlack()でエラー: %w", err))
		DoResponse(ctx, w, errInterServerError, http.StatusInternalServerError)
		return
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(ctx, w, resp, http.StatusOK)
}

// NotifyAccessRankingToSlack : アクセスランキングを Slack に通知
func (n notification) NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyAccessRanking(r.Context())
	if err != nil {
		app.Error(ctx, fmt.Errorf("notificationUseCase.NotifyAccessRanking()でエラー: %w", err))
		DoResponse(ctx, w, errInterServerError, http.StatusInternalServerError)
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(ctx, w, resp, http.StatusOK)
}

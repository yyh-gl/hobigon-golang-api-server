package rest

import (
	"fmt"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Notification : Notification用REST Handlerのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request)
	NotifyAccessRankingToSlack(w http.ResponseWriter, r *http.Request)
	NotifyPokemonEventToSlack(w http.ResponseWriter, r *http.Request)
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
		log.Error(ctx, fmt.Errorf("failed to notificationUseCase.NotifyTodayTasksToSlack(): %w", err))
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
		log.Error(ctx, fmt.Errorf("failed to notificationUseCase.NotifyAccessRanking(): %w", err))
		DoResponse(ctx, w, errInterServerError, http.StatusInternalServerError)
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(ctx, w, resp, http.StatusOK)
}

// NotifyPokemonEventToSlack : Notify event notifications about Pokemon card to Slack.
func (n notification) NotifyPokemonEventToSlack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyPokemonEvent(r.Context())
	if err != nil {
		log.Error(ctx, fmt.Errorf("failed to notificationUseCase.NotifyPokemonEvent(): %w", err))
		DoResponse(ctx, w, errInterServerError, http.StatusInternalServerError)
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(ctx, w, resp, http.StatusOK)
}

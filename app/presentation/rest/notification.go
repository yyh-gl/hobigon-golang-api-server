package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	modelLine "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
	"github.com/yyh-gl/hobigon-golang-api-server/app/log"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Notification : Notification用REST Handlerのインターフェース
type Notification interface {
	NotifyTodayTasksToSlack(w http.ResponseWriter, r *http.Request)
	NotifyPokemonEventToSlack(w http.ResponseWriter, r *http.Request)
	NotifyToLINE(w http.ResponseWriter, r *http.Request)
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

// NotifyPokemonEventToSlack : Notify event notifications about Pokémon card to Slack.
func (n notification) NotifyPokemonEventToSlack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp := notificationResponse{}
	notifiedNum, err := n.u.NotifyPokemonEvent(r.Context())
	if err != nil {
		log.Error(ctx, fmt.Errorf("failed to notificationUseCase.NotifyPokemonEvent(): %w", err))
		DoResponse(ctx, w, errInterServerError, http.StatusInternalServerError)
		return
	}
	resp.NotifiedNum = notifiedNum

	DoResponse(ctx, w, resp, http.StatusOK)
}

// NotifyToLINE : 指定されたBot・メッセージキーでLINEに通知
func (n notification) NotifyToLINE(w http.ResponseWriter, r *http.Request) {
	type request struct {
		BotKey     string `mapstructure:"bot_key" validate:"required"`
		MessageKey string `mapstructure:"message_key" validate:"required"`
	}

	ctx := r.Context()

	var req request
	if err := bindReqWithValidate(ctx, mux.Vars(r), &req); err != nil {
		DoResponse(ctx, w, errBadRequest, http.StatusBadRequest)
		return
	}

	if err := n.u.NotifyToLINE(ctx, req.BotKey, req.MessageKey); err != nil {
		if errors.Is(err, modelLine.ErrLINEMessageKeyNotFound) || errors.Is(err, gateway.ErrLINEBotKeyNotFound) {
			DoResponse(ctx, w, errBadRequest, http.StatusBadRequest)
			return
		}
		log.Error(ctx, fmt.Errorf("failed to notificationUseCase.NotifyToLINE(): %w", err))
		DoResponse(ctx, w, errInterServerError, http.StatusInternalServerError)
		return
	}
	DoResponse(ctx, w, struct{}{}, http.StatusOK)
}

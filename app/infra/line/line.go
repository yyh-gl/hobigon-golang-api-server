package line

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	modelL "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

type line struct {
	channelAccessToken string
}

// NewLine : LINE用のゲートウェイを取得
func NewLine() gateway.Line {
	return &line{
		channelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	}
}

// SendMessage : LINEにメッセージをブロードキャスト配信（LINE公式アカウントの友だち全員に送信）
func (l line) SendMessage(ctx context.Context, msg modelL.Line) error {
	bot, err := messaging_api.NewMessagingApiAPI(l.channelAccessToken)
	if err != nil {
		if app.IsTest() {
			return nil
		}
		return fmt.Errorf("failed to create LINE messaging api client: %w", err)
	}

	_, err = bot.Broadcast(&messaging_api.BroadcastRequest{
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{Text: msg.Text},
		},
	}, "")
	if err != nil && !app.IsTest() {
		return errors.New("failed to send LINE broadcast message: " + err.Error())
	}
	return nil
}

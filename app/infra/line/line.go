package line

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/gateway"
	modelLine "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

// botChannelAccessTokenEnvVars : bot-key -> Channel Access Tokenの環境変数名のホワイトリスト
// 値（トークン自体）は含まないため非機密情報。将来Botを追加する場合はここに1行追加する。
var botChannelAccessTokenEnvVars = map[string]string{
	"son": "SON_LINE_BOT_CHANNEL_ACCESS_TOKEN",
}

type line struct{}

// NewLINE : LINE通知用ゲートウェイを取得
func NewLINE() gateway.LINE {
	return &line{}
}

// SendMessage : 指定されたBotキーに紐づくLINE公式アカウントの友だち全員にメッセージをブロードキャスト配信
func (l line) SendMessage(ctx context.Context, botKey string, msg modelLine.LINE) error {
	envVar, ok := botChannelAccessTokenEnvVars[botKey]
	if !ok {
		return gateway.ErrLINEBotKeyNotFound
	}

	bot, err := messaging_api.NewMessagingApiAPI(os.Getenv(envVar))
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

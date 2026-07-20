package gateway

import (
	"context"
	"errors"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

// ErrLINEBotKeyNotFound : 指定されたBotキーが未知の場合のエラー
var ErrLINEBotKeyNotFound = errors.New("line bot key not found")

// LINE : LINE通知用ゲートウェイのインターフェース
type LINE interface {
	SendMessage(ctx context.Context, botKey string, msg line.LINE) error
}

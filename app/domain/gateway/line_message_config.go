package gateway

import "context"

// LINEMessageConfig : LINE通知メッセージ設定用ゲートウェイのインターフェース
type LINEMessageConfig interface {
	GetMessage(ctx context.Context, messageKey string) (string, error)
}

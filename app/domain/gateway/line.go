package gateway

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/line"
)

// Line : LINE通知用ゲートウェイのインターフェース
type Line interface {
	SendMessage(ctx context.Context, msg line.Line) error
}

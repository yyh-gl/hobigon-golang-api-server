package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/ranking"
)

// Ranking : Ranking用ドメインサービスのインターフェース
type Ranking interface {
	GetAccessRanking(ctx context.Context) (string, ranking.Ranking, error)
}

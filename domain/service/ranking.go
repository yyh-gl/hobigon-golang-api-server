package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/entity"
)

// RankingService : ランキング用サービスのインターフェース
type RankingService interface {
	GetAccessRanking(ctx context.Context) (string, entity.AccessList, error)
}

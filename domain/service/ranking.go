package service

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// RankingService はランキング用サービスのインターフェース
type RankingService interface {
	GetAccessRanking(ctx context.Context) (string, model.AccessList, error)
}

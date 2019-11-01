package service

import (
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// RankingService はランキング用サービスのインターフェース
type RankingService interface {
	GetAccessRanking() (string, model.AccessList, error)
}

package blog

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

// ConvertToDomainModel : DTOからドメインモデルへ変換
func ConvertToDomainModel(ctx context.Context, b dto.BlogDTO) *Blog {
	return &Blog{
		fields{
			Title: Title(b.Title),
			Count: Count(b.Count),
		},
	}
}

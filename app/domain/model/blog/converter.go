package blog

import (
	"context"

	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

// ConvertToEntity : DTOからエンティティへ変換
func ConvertToEntity(ctx context.Context, b dto.BlogDTO) Blog {
	return Blog{
		title: Title(b.Title),
		count: Count(b.Count),
	}
}

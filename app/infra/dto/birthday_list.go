package dto

import (
	"context"
	"fmt"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
)

type BirthdayListDTO []BirthdayDTO

// ConvertToDomainModel : ドメインモデルに変換
func (bl BirthdayListDTO) ConvertToDomainModel(ctx context.Context) (*birthday.BirthdayList, error) {
	var list birthday.BirthdayList
	for _, b := range bl {
		// Birthday モデルを取得
		domainModelBirthday, err := birthday.NewBirthdayWithFullParams(
			b.Name, b.Date, b.WishList,
		)
		if err != nil {
			return nil, fmt.Errorf("birthday.NewBirthdayWithFullParams()でエラー: %w", err)
		}

		list = append(list, *domainModelBirthday)
	}
	return &list, nil
}

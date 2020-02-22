package birthday

import "github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"

// ConvertToDomainModel : DTOからドメインモデルへ変換
func ConvertToDomainModel(b dto.BirthdayDTO) *Birthday {
	return &Birthday{
		fields{
			Name:     Name(b.Name),
			Date:     Date(b.Date),
			WishList: WishList(b.WishList),
		},
	}
}

// ConvertToDomainModelList : リスト型DTOからリスト型ドメインモデルへ変換
func ConvertToDomainModelList(bl dto.BirthdayListDTO) (list BirthdayList) {
	for _, b := range bl {
		model := ConvertToDomainModel(b)
		list = append(list, *model)
	}
	return list
}

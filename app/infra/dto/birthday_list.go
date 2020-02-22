package dto

// BirthdayListDTO : BirthdayDTOのリスト
type BirthdayListDTO []BirthdayDTO

// IsEmpty : BirthdayListDTOが空かどうか判定
func (bl BirthdayListDTO) IsEmpty() bool {
	return len(bl) == 0
}

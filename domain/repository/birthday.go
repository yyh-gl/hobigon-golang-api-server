package repository

import (
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

type BirthdayRepository interface {
	SelectByDate(date string) (model.Birthday, error)
}

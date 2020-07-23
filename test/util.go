package test

import (
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/db"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

func CreateBlog(db *db.DB, title string) {
	_ = db.Create(&dto.BlogDTO{
		Title: title,
		Count: 0,
	}).Error
}

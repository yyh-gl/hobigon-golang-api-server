package rest

import (
	"fmt"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
	model "github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"
	"github.com/yyh-gl/hobigon-golang-api-server/app/usecase"
)

// Birthday : Birthday用REST Handlerのインターフェース
type Birthday interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type birthday struct {
	bu usecase.Birthday
}

// NewBirthday : Birthday用REST Handlerを取得
func NewBirthday(bu usecase.Birthday) Birthday {
	return &birthday{
		bu: bu,
	}
}

// BirthdayResponse : Birthday用共通レスポンス
type BirthdayResponse struct {
	Birthday *model.Birthday `json:"birthday,omitempty"`
	errorResponse
}

// Create : 誕生日データを新規作成
func (b birthday) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string `json:"name" validate:"required,max=30"`
		Date     string `json:"date" validate:"required,len=4"`
		WishList string `json:"wish_list" validate:"required,url"`
	}

	ctx := r.Context()
	req := request{}
	resp := BirthdayResponse{}
	if err := bindReqWithValidate(ctx, &req, r); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusBadRequest)
		return
	}

	birthday, err := b.bu.Create(ctx, req.Name, req.Date, req.WishList)
	if err != nil {
		errInfo := fmt.Errorf("BirthdayUseCase.Create()でエラー: %w", err)
		app.Logger.Println(errInfo)

		resp.Error = errInfo.Error()
		DoResponse(w, resp, http.StatusInternalServerError)
		return
	}
	resp.Birthday = birthday

	DoResponse(w, resp, http.StatusCreated)
}

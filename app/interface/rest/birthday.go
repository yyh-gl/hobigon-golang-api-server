package rest

import (
	"fmt"
	"net/http"
	"time"

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

// birthdayResponse : Birthday用共通レスポンス
// TODO: OK, Error 部分は共通レスポンスにする
type birthdayResponse struct {
	Birthday *model.Birthday `json:"birthday,omitempty"`
}

// Create : 誕生日データを新規作成
func (b birthday) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string    `json:"name" validate:"required,max=30"`
		Date     time.Time `json:"date" validate:"required"`
		WishList string    `json:"wish_list" validate:"required,url"`
	}

	ctx := r.Context()
	req := request{}
	errRes := errorResponse{}
	if err := bindReqWithValidate(ctx, &req, r); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate()でエラー: %w", err)
		app.Logger.Println(errInfo)

		errRes.Error = errInfo.Error()
		DoResponse(w, errRes, http.StatusBadRequest)
		return
	}

	birthday, err := b.bu.Create(ctx, req.Name, req.Date, req.WishList)
	if err != nil {
		errInfo := fmt.Errorf("BirthdayUseCase.Create()でエラー: %w", err)
		app.Logger.Println(errInfo)

		errRes.Error = errInfo.Error()
		DoResponse(w, errRes, http.StatusInternalServerError)
		return
	}

	resp := birthdayResponse{
		Birthday: birthday,
	}
	DoResponse(w, resp, http.StatusCreated)
}

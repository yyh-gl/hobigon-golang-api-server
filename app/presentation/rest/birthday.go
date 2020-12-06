package rest

import (
	"fmt"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
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

// Create : 誕生日データを新規作成
func (b birthday) Create(w http.ResponseWriter, r *http.Request) {
	type (
		request struct {
			Name     string `json:"name" validate:"required,max=30"`
			Date     string `json:"date" validate:"required,len=4"`
			WishList string `json:"wish_list" validate:"required"` // TODO: URL形式のバリデーションを追加
		}
		response struct {
			Name     string `json:"name"`
			Date     string `json:"date"`
			WishList string `json:"wish_list"`
		}
	)

	ctx := r.Context()

	var req request
	if err := bindReqWithValidate(ctx, r, &req); err != nil {
		errInfo := fmt.Errorf("bindReqWithValidate()でエラー: %w", err)
		app.Logger.Println(errInfo)

		DoResponse(w, errBadRequest, http.StatusBadRequest)
		return
	}

	birthday, err := b.bu.Create(ctx, req.Name, req.Date, req.WishList)
	if err != nil {
		errInfo := fmt.Errorf("BirthdayUseCase.Create()でエラー: %w", err)
		app.Logger.Println(errInfo)

		DoResponse(w, errInterServerError, http.StatusInternalServerError)
		return
	}

	resp := response{
		Name:     birthday.Name().String(),
		Date:     birthday.Date().String(),
		WishList: birthday.WishList().String(),
	}
	DoResponse(w, resp, http.StatusCreated)
}

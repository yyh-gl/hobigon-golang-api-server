package rest

import (
	"encoding/json"
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

// TODO: OK, Error 部分は共通レスポンスにする
type birthdayResponse struct {
	OK       bool           `json:"ok"`
	Error    string         `json:"error,omitempty"`
	Birthday model.Birthday `json:"birthday,omitempty"`
}

// Create : 誕生日データを新規作成
func (b birthday) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string    `json:"name"`
		Date     time.Time `json:"date"`
		WishList string    `json:"wish_list"`
	}

	res := birthdayResponse{
		OK: true,
	}

	req, err := decodeRequest(r, request{})
	if err != nil {
		app.Logger.Println(fmt.Errorf("decodeRequest()でエラー: %w", err))

		res.OK = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer func() { _ = r.Body.Close() }()

	if res.OK {
		// リクエストパラメータ内の date を time.Time 型に変換
		// TODO: フォーマット部分を定数化
		date, err := time.Parse("2006-01-02T15:04:05.000000Z", req["date"].(string))
		if err != nil {
			app.Logger.Println(fmt.Errorf("time.Parse()内でエラー: %w", err))

			res.OK = false
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		}

		if res.OK {
			createdBirthday, err := b.bu.Create(r.Context(), req["name"].(string), date, req["wish_list"].(string))
			if err != nil {
				app.Logger.Println(fmt.Errorf("BirthdayUseCase.Create()でエラー: %w", err))

				res.OK = false
				res.Error = err.Error()
				w.WriteHeader(http.StatusInternalServerError)
			}
			res.Birthday = *createdBirthday
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		app.Logger.Println(fmt.Errorf("json.NewEncoder().Encode()でエラー: %w", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

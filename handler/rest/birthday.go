package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model/birthday"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/usecase"
)

//////////////////////////////////////////////////
// NewBirthdayHandler
//////////////////////////////////////////////////

// BirthdayHandler : ブログ用のハンドラーインターフェース
type BirthdayHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type birthdayHandler struct {
	bu usecase.BirthdayUseCase
}

// NewBirthdayHandler : ブログ用のハンドラーを取得
func NewBirthdayHandler(bu usecase.BirthdayUseCase) BirthdayHandler {
	return &birthdayHandler{
		bu: bu,
	}
}

// TODO: OK, Error 部分は共通レスポンスにする
type birthdayResponse struct {
	OK       bool                  `json:"ok"`
	Error    string                `json:"error,omitempty"`
	Birthday birthday.BirthdayJSON `json:"birthday,omitempty"`
}

//////////////////////////////////////////////////
// Create
//////////////////////////////////////////////////

// Create : 誕生日データを新規作成
func (bh birthdayHandler) Create(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Name     string    `json:"name"`
		Date     time.Time `json:"date"`
		WishList string    `json:"wish_list"`
	}

	logger := app.Logger

	res := birthdayResponse{
		OK: true,
	}

	req, err := decodeRequest(r, request{})
	if err != nil {
		logger.Println(err)

		res.OK = false
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if res.OK {
		// リクエストパラメータ内の date を time.Time 型に変換
		// TODO: フォーマット部分を定数化
		date, err := time.Parse("2006-01-02T15:04:05.000000Z", req["date"].(string))
		if err != nil {
			e := errors.Wrap(err, "time.Parse()内でのエラー")
			logger.Println(e)

			res.OK = false
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
		}

		if res.OK {
			birthday, err := bh.bu.Create(r.Context(), req["name"].(string), date, req["wish_list"].(string))
			if err != nil {
				logger.Println(err)

				res.OK = false
				res.Error = err.Error()
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				// JSON 形式に変換
				res.Birthday = birthday.JSONSerialize()
			}
		}
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

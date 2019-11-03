package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
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

type birthday struct {
	ID        uint       `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Date      string     `json:"date,omitempty"`
	WishList  string     `json:"wish_list,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TODO: OK, Error 部分は共通レスポンスにする
type birthdayResponse struct {
	OK    bool      `json:"ok"`
	Error string    `json:"error,omitempty"`
	Blog  *birthday `json:"birthday,omitempty"`
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

	var b *model.Birthday
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
			b, err = bh.bu.Create(r.Context(), req["name"].(string), date, req["wish_list"].(string))
			if err != nil {
				logger.Println(err)

				res.OK = false
				res.Error = err.Error()
				w.WriteHeader(http.StatusInternalServerError)
			}

			// Birthday モデルをレスポン用に変換
			if b != nil {
				birthdayRes := birthday{
					ID:        b.ID(),
					Name:      b.Name(),
					Date:      b.Date().String(),
					WishList:  b.WishList().String(),
					CreatedAt: b.CreatedAt(),
					UpdatedAt: b.UpdatedAt(),
					DeletedAt: b.DeletedAt(),
				}
				res.Blog = &birthdayRes
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

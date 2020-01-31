package rest_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/app/domain/model/birthday"

	"github.com/bmizerany/assert"
	"github.com/yyh-gl/hobigon-golang-api-server/app/interface/rest"
	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

func createBirthday(c *test.Client, name string, date string, wishList string) {
	c.AddRoute(http.MethodPost, "/api/v1/birthday", c.DIContainer.HandlerBirthday.Create)

	body := `{
	"name": "` + name + `",
	"date": "` + date + `",
	"wish_list": "` + wishList + `"
}`
	_ = c.Post("/api/v1/birthday", body)
}

func TestBirthdayHandler_Create(t *testing.T) {
	c := test.NewClient()
	defer func() { _ = c.DIContainer.DB.Close() }()

	c.AddRoute(http.MethodPost, "/api/v1/birthday", c.DIContainer.HandlerBirthday.Create)

	testCases := []struct {
		name         string
		date         string
		wishList     string
		wantName     string
		wantDate     string
		wantWishList string
		err          string
	}{
		{ // 正常系
			name:         "hon-D",
			date:         "1205",
			wishList:     "https://honzon.co.jp",
			wantName:     "hon-D",
			wantDate:     "1205",
			wantWishList: "https://honzon.co.jp",
			err:          "",
		},
		{ // 正常系：重複はOK
			name:         "duplicate-name",
			date:         "1205",
			wishList:     "https://honzon.co.jp",
			wantName:     "duplicate-name",
			wantDate:     "1205",
			wantWishList: "https://honzon.co.jp",
			err:          "",
		},
		{ // 正常系：30文字name
			name:         "hon-Dhon-Dhon-Dhon-Dhon-Dhon-D",
			date:         "1205",
			wishList:     "https://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "",
		},
		{ // 異常系：nameがない
			name:         "",
			date:         "1205",
			wishList:     "https://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{ // 異常系：31文字name
			name:         "hon-Dhon-Dhon-Dhon-Dhon-Dhon-D1",
			date:         "1205",
			wishList:     "https://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Name' Error:Field validation for 'Name' failed on the 'max' tag",
		},
		{ // 異常系：dateがない
			name:         "hon-D",
			date:         "",
			wishList:     "https://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Date' Error:Field validation for 'Date' failed on the 'required' tag",
		},
		{ // 異常系：3文字date
			name:         "hon-D",
			date:         "125",
			wishList:     "https://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Date' Error:Field validation for 'Date' failed on the 'len' tag",
		},
		{ // 異常系：5文字date
			name:         "hon-D",
			date:         "11205",
			wishList:     "https://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.Date' Error:Field validation for 'Date' failed on the 'len' tag",
		},
		{ // 異常系：wish_listがない
			name:         "hon-D",
			date:         "1205",
			wishList:     "",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "bindReqWithValidate()でエラー: バリデーションエラー: Key: 'request.WishList' Error:Field validation for 'WishList' failed on the 'required' tag",
		},
		{ // 異常系：wish_listに不正な形式のURLが入っている
			name:         "hon-D",
			date:         "1205",
			wishList:     "httpsa://honzon.co.jp",
			wantName:     "",
			wantDate:     "",
			wantWishList: "",
			err:          "BirthdayUseCase.Create()でエラー: model.NewBirthday()内でエラー: NewWishList()内でエラー: バリデーションエラー：【Birthday】WishListが\"https://\"から始まっていません",
		},
	}

	for _, tc := range testCases {
		body := `{
	"name": "` + tc.name + `",
	"date": "` + tc.date + `",
	"wish_list": "` + tc.wishList + `"
}`
		rec := c.Post("/api/v1/birthday", body)
		resp := rest.BirthdayResponse{}
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)

		// 重複データ登録時に使用するテストデータを追加
		createBirthday(c, "duplicate-name", "1205", "https://honzon.co.jp")

		if tc.err == "" {
			if tc.wantName != "" {
				assert.Equal(t, tc.wantName, resp.Birthday.Name().String())
			}
			if tc.wantDate != "" {
				assert.Equal(t, tc.wantDate, resp.Birthday.Date().String())
			}
			if tc.wantWishList != "" {
				assert.Equal(t, tc.wantWishList, resp.Birthday.WishList().String())
			}
			assert.Equal(t, "", resp.Error)
		} else {
			assert.Equal(t, (*birthday.Birthday)(nil), resp.Birthday)
			assert.Equal(t, tc.err, resp.Error)
		}
	}
}

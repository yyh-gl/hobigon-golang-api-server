package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yyh-gl/hobigon-golang-api-server/test"
)

func TestBirthdayHandler_Create(t *testing.T) {
	type want struct {
		body       string
		statusCode int
	}

	tests := []struct {
		name     string
		date     string
		wishList string
		want     want
	}{
		{ // 正常系
			name:     "hon-D",
			date:     "0905",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"name":"hon-D","date":"0905","wish_list":"https://honzon.co.jp"}`,
				statusCode: http.StatusCreated,
			},
		},
		{ // 正常系：重複はOK
			name:     "duplicate-name",
			date:     "0905",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"name":"duplicate-name","date":"0905","wish_list":"https://honzon.co.jp"}`,
				statusCode: http.StatusCreated,
			},
		},
		{ // 正常系：30文字name
			name:     "hon-Dhon-Dhon-Dhon-Dhon-Dhon-D",
			date:     "0905",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"name":"hon-Dhon-Dhon-Dhon-Dhon-Dhon-D","date":"0905","wish_list":"https://honzon.co.jp"}`,
				statusCode: http.StatusCreated,
			},
		},
		{ // 異常系：nameがない
			name:     "",
			date:     "0905",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：31文字name
			name:     "hon-Dhon-Dhon-Dhon-Dhon-Dhon-D1",
			date:     "0905",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：dateがない
			name:     "hon-D",
			date:     "",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：3文字date
			name:     "hon-D",
			date:     "125",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：5文字date
			name:     "hon-D",
			date:     "10905",
			wishList: "https://honzon.co.jp",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：wish_listがない
			name:     "hon-D",
			date:     "0905",
			wishList: "",
			want: want{
				body:       `{"error":{"detail":"不正なリクエスト形式です"}}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{ // 異常系：wish_listに不正な形式のURLが入っている
			name:     "hon-D",
			date:     "0905",
			wishList: "httpsa://honzon.co.jp",
			// FIXME: BadRequestエラーを返す
			want: want{
				body:       `{"error":{"detail":"サーバ内でエラーが発生しました"}}`,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	// 重複データ登録時に使用するテストデータを追加
	test.CreateBirthday(DIContainer.DB, "duplicate-name", "0905", "https://honzon.co.jp")

	for _, tt := range tests {
		reqBody := strings.NewReader(`
{
	"name": "` + tt.name + `",
	"date": "` + tt.date + `",
	"wish_list": "` + tt.wishList + `"
}
`)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/birthday", reqBody)
		rr := httptest.NewRecorder()
		Router.ServeHTTP(rr, req)

		if c := rr.Code; c != tt.want.statusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				c, tt.want.statusCode)
		}

		respBody := strings.TrimRight(rr.Body.String(), "\n")
		if respBody != tt.want.body {
			t.Errorf("handler returned unexpected body: got %v want %v",
				respBody, tt.want.body)
		}
	}
}

package handler

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

type NotifyBirthdayRequest struct {
	Date string `json:"date"`
}

func NotifyBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ctx.Value("logger").(*log.Logger)

	birthdayRepository := repository.NewBirthdayRepository()
	slackGateway := gateway.NewSlackGateway()

	// TODO: デコード処理を共通化
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var notifyBirthdayRequest NotifyBirthdayRequest
	err = json.Unmarshal(body, &notifyBirthdayRequest)
	if err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	birthday, err := birthdayRepository.SelectByDate(ctx, notifyBirthdayRequest.Date)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if birthday.IsToday() {
		err = slackGateway.SendBirthday(ctx, birthday)
		if err != nil {
			logger.Println(err)
			// TODO: エラーハンドリングをきちんとする
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(birthday); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

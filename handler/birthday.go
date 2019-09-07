package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

type NotifyBirthdayRequest struct {
	Date string `json:"date"`
}

func NotifyBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	ctx := r.Context()

	birthdayRepository := repository.NewBirthdayRepository()
	slackGateway := gateway.NewSlackGateway()

	// TODO: body から受け取る必要なくない？
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	birthday, err := birthdayRepository.SelectByDate(notifyBirthdayRequest.Date)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if birthday.IsToday() {
		err = slackGateway.SendBirthday(ctx, birthday)
		if err != nil {
			logger.Println(err)
			// TODO: エラーハンドリングをきちんとする
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(birthday); err != nil {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

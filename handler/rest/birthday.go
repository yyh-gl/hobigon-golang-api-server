package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/repository"
)

func NotifyBirthdayHandler(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger

	birthdayRepository := repository.NewBirthdayRepository()
	slackGateway := gateway.NewSlackGateway()

	today := time.Now().Format("0102")
	birthday, err := birthdayRepository.SelectByDate(today)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Println(err)
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if birthday.IsToday() {
		err = slackGateway.SendBirthday(birthday)
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

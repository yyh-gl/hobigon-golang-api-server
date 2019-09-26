package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/infra"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// TODO: robo タスクとしても実行できるようにしたい
func GetAccessRanking(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Success bool `json:"success"`
	}

	logger := app.Logger

	//slackGateway := gateway.NewSlackGateway()

	hoge, _ := infra.GetAccessRanking()
	fmt.Println("========================")
	fmt.Println(hoge)
	fmt.Println("========================")

	responseJSON, err := json.Marshal(response{Success: true})
	if err != nil {
		logger.Println(err)
		http.Error(w, "Error at json.Marshal()", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseJSON)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Error at w.Write()", http.StatusInternalServerError)
		return
	}
}

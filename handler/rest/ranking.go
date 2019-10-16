package rest

import (
	"encoding/json"
	"net/http"

	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"

	"github.com/yyh-gl/hobigon-golang-api-server/infra/gateway"

	"github.com/yyh-gl/hobigon-golang-api-server/infra"

	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

// TODO: robo タスクとしても実行できるようにしたい
func GetAccessRanking(w http.ResponseWriter, r *http.Request) {
	// FIXME: もし他にもランキングができてきたら、このレスポンス用構造体を関数外に出して共通化することで、
	//  レスポンスの形を統一できるから、解析処理とかクライアント側の処理が捗りそう
	type response struct {
		Ranking model.AccessList `json:"ranking"`
	}

	logger := app.Logger

	slackGateway := gateway.NewSlackGateway()

	// アクセスランキングの結果を取得
	// TODO: エクセルに出力して解析とかしたい
	rankingMsg, accessList, err := infra.GetAccessRanking()
	if err != nil {
		logger.Println(err)
		http.Error(w, "Error at infra.GetAccessRanking()", http.StatusInternalServerError)
		return
	}

	// アクセスランキングの結果を Slack に通知
	err = slackGateway.SendRanking(rankingMsg)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Error at slackGateway.SendRanking()", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(response{Ranking: accessList})
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

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type response struct {
	Success string
}

func TaskHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	response, _ := json.Marshal(response{ Success: "ok"})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// TODO: エラーハンドリングをきちんとする
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

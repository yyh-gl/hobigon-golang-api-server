package handler

import (
	"encoding/json"
	"net/http"
)

// ============================================
// 技術検証用ハンドラー
// ============================================

// ヘッダーに関するHTML（文字列）を返すハンドラー
func GetHeaderHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Header string `json:"header"`
	}

	res := response{
		Header: `
<header class="App-header">
    <h1>ヘッダー</h1>
</header>
`,
	}

	//	time.Sleep(1 * time.Second)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// フッターに関するHTML（文字列）を返すハンドラー
func GetFooterHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Footer string `json:"footer"`
	}

	res := response{
		Footer: `
<footer class="App-footer">
    <h1>フッター</h1>
</footer>
`,
	}

	//	time.Sleep(1 * time.Second)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

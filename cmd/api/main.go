package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
)

func main() {
	// 依存関係を定義
	diContainer := initApp()
	defer func() { _ = diContainer.DB.Close() }()

	// ロガー設定
	// TODO: いちいちdi.Containerにバインドする意味があるのかもう一度検討
	app.Logger = diContainer.Logger

	// ルーティング設定
	r := httprouter.New()
	r.GlobalOPTIONS = wrapHandler(preflightHandler)

	// ブログ関連のAPI
	r.HandlerFunc(http.MethodPost, "/api/v1/blogs", wrapHandler(diContainer.HandlerBlog.Create))
	r.HandlerFunc(http.MethodGet, "/api/v1/blogs/:title", wrapHandler(diContainer.HandlerBlog.Show))
	r.HandlerFunc(http.MethodPost, "/api/v1/blogs/:title/like", wrapHandler(diContainer.HandlerBlog.Like))

	// 誕生日関連のAPI
	r.HandlerFunc(http.MethodPost, "/api/v1/birthday", wrapHandler(diContainer.HandlerBirthday.Create))

	// 通知系API
	r.HandlerFunc(http.MethodPost, "/api/v1/notifications/slack/tasks/today", wrapHandler(diContainer.HandlerNotification.NotifyTodayTasksToSlack))
	// TODO: 誕生日の人が複数いたときに対応
	r.HandlerFunc(http.MethodPost, "/api/v1/notifications/slack/birthdays/today", wrapHandler(diContainer.HandlerNotification.NotifyTodayBirthdayToSlack))
	r.HandlerFunc(http.MethodPost, "/api/v1/notifications/slack/rankings/access", wrapHandler(diContainer.HandlerNotification.NotifyAccessRankingToSlack))

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println(" ↳  Log File -> " + os.Getenv("LOG_PATH") + "/" + app.APILogFilename)
	fmt.Println("========================")
	app.Logger.Println("Server Start")
	app.Logger.Fatal(http.ListenAndServe(":3000", r))
}

// wrapHandler : 全ハンドラー共通処理
func wrapHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエスト内容をログ出力
		// TODO: Bodyの内容を出力
		app.Logger.Println("[AccessLog] " + r.Method + " " + r.URL.String())

		// CORS用ヘッダーを付与
		switch {
		case app.IsPrd():
			w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
		case app.IsDev() || app.IsTest():
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
		}
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		h.ServeHTTP(w, r)
	}
}

// preflightHandler : preflight用のハンドラー
func preflightHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

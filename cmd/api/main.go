package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/context"
)

func main() {
	// システム共通で使用するものを用意
	//  -> logger, DB
	app.Init(app.APILogFilename)
	defer app.DB.Close()

	// 依存関係を定義
	blogHandler := initBlogHandler()
	birthdayHandler := initBirthdayHandler()
	notificationHandler := initNotificationHandler()

	// ルーティング設定
	r := httprouter.New()
	r.GlobalOPTIONS = http.HandlerFunc(corsHandler)
	//r.OPTIONS("/*path", corsHandler) // CORS用の pre-flight 設定

	// ブログ関連のAPI
	r.POST("/api/v1/blogs", wrapHandler(http.HandlerFunc(blogHandler.Create)))
	r.GET("/api/v1/blogs/:title", wrapHandler(http.HandlerFunc(blogHandler.Show)))
	r.POST("/api/v1/blogs/:title/like", wrapHandler(http.HandlerFunc(blogHandler.Like)))

	// 誕生日関連のAPI
	r.POST("/api/v1/birthday", wrapHandler(http.HandlerFunc(birthdayHandler.Create)))

	// 通知系API
	r.POST("/api/v1/notifications/slack/tasks/today", wrapHandler(http.HandlerFunc(notificationHandler.NotifyTodayTasksToSlack)))
	// TODO: 誕生日の人が複数いたときに対応
	r.POST("/api/v1/notifications/slack/birthdays/today", wrapHandler(http.HandlerFunc(notificationHandler.NotifyTodayBirthdayToSlack)))
	r.POST("/api/v1/notifications/slack/rankings/access", wrapHandler(http.HandlerFunc(notificationHandler.NotifyAccessRankingToSlack)))

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println(" ↳  Log File -> " + os.Getenv("LOG_PATH") + "/" + app.APILogFilename)
	fmt.Println("========================")
	app.Logger.Print("Server Start")
	app.Logger.Fatal(http.ListenAndServe(":3000", r))
}

func corsHandler(w http.ResponseWriter, _ *http.Request) {
	switch {
	case app.IsPrd():
		w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
	case app.IsDev() || app.IsTest():
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
	}
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.InjectRequestParams(ctx, ps)
		r = r.WithContext(ctx)

		// リクエスト内容をログ出力
		// TODO: Body の内容を記録
		app.Logger.Print("[AccessLog] " + r.Method + " " + r.URL.String())

		// 共通ヘッダー設定
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

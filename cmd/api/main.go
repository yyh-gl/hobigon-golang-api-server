package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/context"

	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/handler"
)

func main() {
	// システム共通で使用するものを用意
	//  -> logger, DB
	app.Init()

	// ルーティング設定
	r := httprouter.New()
	r.OPTIONS("/*path", corsHandler)                                                  // CORS用の pre-flight 設定
	r.POST("/api/v1/tasks", wrapHandler(http.HandlerFunc(handler.NotifyTaskHandler))) // Slack 通知のために POST メソッド
	r.POST("/api/v1/blogs", wrapHandler(http.HandlerFunc(handler.CreateBlogHandler)))
	r.GET("/api/v1/blogs/:title", wrapHandler(http.HandlerFunc(handler.GetBlogHandler)))
	r.POST("/api/v1/blogs/:title/like", wrapHandler(http.HandlerFunc(handler.LikeBlogHandler)))
	r.POST("/api/v1/birthdays/today", wrapHandler(http.HandlerFunc(handler.NotifyBirthdayHandler))) // Slack 通知のために POST メソッド
	r.POST("/api/v1/rankings/access", wrapHandler(http.HandlerFunc(handler.GetAccessRanking)))      // Slack 通知のために POST メソッド

	// 技術検証用ルーティング設定
	//r.GET("/api/v1/header", wrapHandler(http.HandlerFunc(handler.GetHeaderHandler)))
	//r.GET("/api/v1/footer", wrapHandler(http.HandlerFunc(handler.GetFooterHandler)))

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println(" ↳  Log File -> " + os.Getenv("LOG_PATH") + "/app.log")
	fmt.Println("========================")
	app.Logger.Print("Server Start")
	app.Logger.Fatal(http.ListenAndServe(":3000", r))
}

func corsHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
	//w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
	//w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
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
		w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
		//w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
		//w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3001")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		h.ServeHTTP(w, r)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
	"github.com/yyh-gl/hobigon-golang-api-server/handler"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// ロガー設定
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// ログ出力先を設定
	logPath := os.Getenv("LOG_PATH")
	logfile, err := os.OpenFile(logPath + "/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open " + logPath + "/app.log:" + err.Error())
	}
	defer logfile.Close()

	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		logger.SetOutput(logfile)
	default:
		logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	}

	// ルーティング設定
	r := httprouter.New()
	r.POST("/api/v1/tasks", wrapHandler(http.HandlerFunc(handler.NotifyTaskHandler), *logger))
	r.POST("/api/v1/blogs", wrapHandler(http.HandlerFunc(handler.CreateBlogHandler), *logger))
	r.GET("/api/v1/blogs", wrapHandler(http.HandlerFunc(handler.GetBlogHandler), *logger))
	r.POST("/api/v1/blogs/like", wrapHandler(http.HandlerFunc(handler.LikeBlogHandler), *logger))

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println("========================")
	logger.Fatal(http.ListenAndServe(":3000", r))
}

func wrapHandler(h http.Handler, logger log.Logger) httprouter.Handle {
	// DB設定
	db := getGormConnect()
	// TODO: 動作チェック
	//defer db.Close()

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "params", ps)
		ctx = context.WithValue(ctx, "logger", logger)
		ctx = context.WithValue(ctx, "db", db)
		r = r.WithContext(ctx)

		// 共通ヘッダー設定
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func getGormConnect() *gorm.DB {
	DBMS     := "mysql"
	USER     := os.Getenv("MYSQL_USER")
	PASSWORD := os.Getenv("MYSQL_PASSWORD")
	PROTOCOL := "tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")"
	DATABASE := os.Getenv("MYSQL_DATABASE")

	// ?parseTime=true によりレコードSELECT時のスキャンエラーとやらを無視できる
	CONNECT := USER+":"+PASSWORD+"@"+PROTOCOL+"/"+DATABASE+"?parseTime=true&loc=Asia%2FTokyo"

	db,err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}

	// マイグレーション実行
	db.AutoMigrate(&model.Blog{})

	return db
}

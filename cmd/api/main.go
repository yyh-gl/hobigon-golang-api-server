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

type server struct {
	router *httprouter.Router
	logger *log.Logger
	db     *gorm.DB
}

func initDependencies() server {
	return server{
		router: httprouter.New(),
		logger: getLogger(),
		db:     getGormConnect(),
	}
}

func main() {
	s := initDependencies()
	defer s.db.Close()

	// ルーティング設定
	s.router.OPTIONS("/*path", corsHandler) // CORS用の pre-flight 設定
	s.router.POST("/api/v1/tasks", wrapHandler(http.HandlerFunc(handler.NotifyTaskHandler), s))
	s.router.POST("/api/v1/blogs", wrapHandler(http.HandlerFunc(handler.CreateBlogHandler), s))
	s.router.GET("/api/v1/blogs/:title", wrapHandler(http.HandlerFunc(handler.GetBlogHandler), s))
	s.router.POST("/api/v1/blogs/:title/like", wrapHandler(http.HandlerFunc(handler.LikeBlogHandler), s))
	s.router.POST("/api/v1/birthdays/today", wrapHandler(http.HandlerFunc(handler.NotifyBirthdayHandler), s))

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println(" ↳  Log File -> " + os.Getenv("LOG_PATH") + "/app.log")
	fmt.Println("========================")
	s.logger.Print("Server Start")
	s.logger.Fatal(http.ListenAndServe(":3000", s.router))
}

func corsHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	//w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:1313")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func wrapHandler(h http.Handler, s server) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "params", ps)
		ctx = context.WithValue(ctx, "logger", *s.logger)
		ctx = context.WithValue(ctx, "db", s.db)
		r = r.WithContext(ctx)

		// 共通ヘッダー設定
		w.Header().Add("Access-Control-Allow-Origin", "https://yyh-gl.github.io")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")

		h.ServeHTTP(w, r)
	}
}

func getLogger() *log.Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// ログ出力先を設定
	logPath := os.Getenv("LOG_PATH")
	logfile, err := os.OpenFile(logPath + "/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open " + logPath + "/app.log:" + err.Error())
	}

	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		logger.SetOutput(logfile)
	default:
		logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	}

	return logger
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
	db.AutoMigrate(&model.Birthday{})

	return db
}

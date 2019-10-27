package app

import (
	"io"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// justifying
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yyh-gl/hobigon-golang-api-server/domain/model"
)

// TODO: api と cli で分ける（それぞれの main の中に入れてしまってもいいかも）

var (
	Logger *log.Logger
	DB     *gorm.DB
)

// コンテキストにセットするさいのキー用の型
type contextKey int

// CliContextKey は cli.Context を context.Context にセットするさいのキー
const CliContextKey contextKey = iota

// ログファイル名
const (
	APILogFilename string = "api.log"
	CLiLogFilename string = "cli.log"
)

// アプリ全体で使用する機能を初期化する関数
func Init(logFilename string) {
	Logger = getLogger(logFilename)
	DB = getGormConnect()
}

// getLogger はロガーを取得する関数
func getLogger(filename string) *log.Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// ログ出力先を設定
	logPath := os.Getenv("LOG_PATH")
	logfile, err := os.OpenFile(logPath+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open " + logPath + "/" + filename + err.Error())
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

// getGormConnect は DB コネクションを取得する関数
func getGormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := os.Getenv("MYSQL_USER")
	PASSWORD := os.Getenv("MYSQL_PASSWORD")
	PROTOCOL := "tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")"
	DATABASE := os.Getenv("MYSQL_DATABASE")

	// ?parseTime=true によりレコードSELECT時のスキャンエラーとやらを無視できる
	CONNECT := USER + ":" + PASSWORD + "@" + PROTOCOL + "/" + DATABASE + "?parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}

	// マイグレーション実行
	db.AutoMigrate(&model.Blog{})
	db.AutoMigrate(&model.Birthday{})

	return db
}

// IsPrd は実行環境が Production か確認する関数
func IsPrd() bool {
	return os.Getenv("APP_ENV") == "production"
}

package app

import (
	"io"
	"log"
	"os"

	// justifying
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: api と cli で分ける（それぞれの main の中に入れてしまってもいいかも）

// Logger : システム共通ロガー
var Logger *log.Logger

// コンテキストにセットするさいのキー用の型
type contextKey int

const (
	// CliContextKey : cli.Context を context.Context にセットするさいのキー
	CliContextKey contextKey = iota

	// APILogFilename : APIサーバ関連のログファイル名
	APILogFilename string = "api.log"
	// CLILogFilename : CLI関連のログファイル名
	CLILogFilename string = "cli.log"
)

// NewLogger : ロガーを生成
func NewLogger(filename string) *log.Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	if IsTest() {
		logger.SetOutput(os.Stdout)
		return logger
	}

	// ログ出力先を設定
	logPath := os.Getenv("LOG_PATH")
	logfile, err := os.OpenFile(logPath+"/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open " + logPath + "/" + filename + ": " + err.Error())
	}

	if IsPrd() {
		logger.SetOutput(logfile)
	} else {
		logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	}

	return logger
}

// IsDev : 実行環境が Development かどうかを確認
func IsDev() bool {
	return os.Getenv("APP_ENV") == "develop"
}

// IsTest : 実行環境が Test かどうかを確認
func IsTest() bool {
	return os.Getenv("APP_ENV") == "test"
}

// IsPrd : 実行環境が Production かどうかを確認
func IsPrd() bool {
	return os.Getenv("APP_ENV") == "production"
}

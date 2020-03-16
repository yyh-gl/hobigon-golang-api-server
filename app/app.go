package app

import (
	"log"
	"os"
)

// Logger : システム共通ロガー
var Logger *log.Logger

// contextKey : コンテキストにセットするさいのキー用の型
type contextKey int

// CLIContextKey : cli.Contextをcontext.Contextにセットするさいのキー
const CLIContextKey contextKey = iota

// IsDev : 実行環境がDevelopmentかどうかを確認
func IsDev() bool {
	return os.Getenv("APP_ENV") == "develop"
}

// IsTest : 実行環境がTestかどうかを確認
func IsTest() bool {
	return os.Getenv("APP_ENV") == "test"
}

// IsPrd : 実行環境がProductionかどうかを確認
func IsPrd() bool {
	return os.Getenv("APP_ENV") == "production"
}

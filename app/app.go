package app

import (
	"log"
	"os"

	// justifying
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// Logger : システム共通ロガー
var Logger *log.Logger

// contextKey : コンテキストにセットするさいのキー用の型
type contextKey int

// CLIContextKey : cli.Context を context.Context にセットするさいのキー
const CLIContextKey contextKey = iota

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

package app

import (
	"io"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// justifying
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/imodel"
)

// TODO: api と cli で分ける（それぞれの main の中に入れてしまってもいいかも）

var (
	// Logger : システム共通ロガー
	Logger *log.Logger
	// DB : システム共通 DB クライアント
	DB *gorm.DB
)

// コンテキストにセットするさいのキー用の型
type contextKey int

const (
	// CliContextKey : cli.Context を context.Context にセットするさいのキー
	CliContextKey contextKey = iota

	// APILogFilename : APIサーバ関連のログファイル名
	APILogFilename string = "api.log"
	// CLILogFilename : CLI関連のログファイル名
	CLILogFilename string = "cli.log"

	// SQLiteDBFile : SQLite3のローカル用DBデータ
	SQLiteDBFile string = "local.db"
)

// Init : アプリ全体で使用する機能を初期化
func Init(logFilename string) {
	Logger = newLogger(logFilename)
	if IsPrd() {
		DB = newMySQLConnect()
	} else {
		DB = newSQLiteConnect()
	}
}

// newLogger : ロガーを生成
func newLogger(filename string) *log.Logger {
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

// newMySQLConnect : DB（MySQL）コネクションを生成
func newMySQLConnect() *gorm.DB {
	dbms := "mysql"
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	protocol := "tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")"
	database := os.Getenv("MYSQL_DATABASE")

	// charset=utf8mb4 により文字コードを utf8mb4 に変更
	// parseTime=true によりレコードSELECT時のスキャンエラーとやらを無視できる
	CONNECT := user + ":" + password + "@" + protocol + "/" + database + "?charset=utf8mb4,utf8&parseTime=true&loc=Asia%2FTokyo"

	db, err := gorm.Open(dbms, CONNECT)
	if err != nil {
		panic(err.Error())
	}

	migrate(db)

	return db
}

// newSQLiteConnect : DB（SQLite）コネクションを生成
func newSQLiteConnect() (db *gorm.DB) {
	dbms := "sqlite3"

	var err error
	if IsDev() {
		db, err = gorm.Open(dbms, SQLiteDBFile)
	} else if IsTest() {
		db, err = gorm.Open("sqlite3", ":memory:")
	}
	if err != nil {
		panic(err.Error())
	}

	migrate(db)

	return db
}

// migrate : マイグレーション実施
func migrate(db *gorm.DB) {
	db.AutoMigrate(&imodel.BlogDTO{}, &imodel.BirthdayDTO{})
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

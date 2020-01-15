package db

import (
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/app"

	"github.com/jinzhu/gorm"
	"github.com/yyh-gl/hobigon-golang-api-server/infra/imodel"
)

// sqliteDBFile : SQLite3のローカル用DBデータ
const sqliteDBFile string = "local.db"

// DB : DBコネクション->infra以外でORMライブラリを意識させないための中間層
type DB = gorm.DB

// NewDB : DBコネクションを生成
func NewDB() (db *DB) {
	if app.IsTest() {
		db = newSQLiteConnect()
	} else {
		db = newMySQLConnect()
	}
	return db
}

// newMySQLConnect : DB（MySQL）コネクションを生成
func newMySQLConnect() *DB {
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
func newSQLiteConnect() *DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err.Error())
	}

	migrate(db)

	return db
}

// migrate : マイグレーション実施
func migrate(db *DB) {
	db.AutoMigrate(&imodel.BlogDTO{}, &imodel.BirthdayDTO{})
}

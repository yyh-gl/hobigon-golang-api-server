package db

import (
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yyh-gl/hobigon-golang-api-server/app"
	"github.com/yyh-gl/hobigon-golang-api-server/app/infra/dto"
)

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
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err.Error())
	}

	conf := mysql.Config{
		Net:       "tcp",
		Addr:      os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		DBName:    os.Getenv("MYSQL_DATABASE"),
		User:      os.Getenv("MYSQL_USER"),
		Passwd:    os.Getenv("MYSQL_PASSWORD"),
		Loc:       jst,
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
	}

CONNECT:

	db, err := gorm.Open("mysql", conf.FormatDSN())
	if err != nil {
		if app.IsDev() {
			if strings.Contains(err.Error(), "connection refused") {
				time.Sleep(time.Second)
				goto CONNECT
			}
		}

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
	db.AutoMigrate(&dto.BlogDTO{})
}

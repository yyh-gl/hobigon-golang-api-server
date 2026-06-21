package db

import (
	"os"
	"strings"
	"time"

	mysqld "github.com/go-sql-driver/mysql"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

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

	conf := mysqld.Config{
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

	db, err := gorm.Open(mysql.Open(conf.FormatDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		if app.IsDev() {
			if strings.Contains(err.Error(), "connection refused") {
				time.Sleep(time.Second)
				goto CONNECT
			}
		}

		panic(err.Error())
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err.Error())
	}

	migrate(db)

	return db
}

// newSQLiteConnect : DB（SQLite）コネクションを生成
func newSQLiteConnect() *DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err.Error())
	}

	migrate(db)

	return db
}

// migrate : マイグレーション実施
func migrate(db *DB) {
	if err := db.AutoMigrate(&dto.BlogDTO{}); err != nil {
		panic(err.Error())
	}
}

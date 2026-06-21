package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	mysqld "github.com/go-sql-driver/mysql"
	"github.com/glebarez/sqlite"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

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

	registerOTelCallbacks(db)
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

const otelSpanKey = "otel:span"

func registerOTelCallbacks(db *DB) {
	tracer := otel.Tracer("hobigon-gorm")

	beforeFn := func(opName string) func(*gorm.DB) {
		return func(tx *gorm.DB) {
			if tx.Statement == nil || tx.Statement.Context == nil {
				return
			}
			ctx, span := tracer.Start(tx.Statement.Context, fmt.Sprintf("gorm:%s", opName),
				trace.WithAttributes(
					attribute.String(string(semconv.DBSystemKey), "mysql"),
				),
			)
			tx.Statement.Context = ctx
			tx.Statement.InstanceSet(otelSpanKey, span)
		}
	}

	endSpan := func(tx *gorm.DB) {
		if tx.Statement == nil {
			return
		}
		val, ok := tx.Statement.InstanceGet(otelSpanKey)
		if !ok {
			return
		}
		span, ok := val.(trace.Span)
		if !ok {
			return
		}
		defer span.End()
		if tx.Error != nil {
			span.RecordError(tx.Error)
			span.SetStatus(codes.Error, tx.Error.Error())
		}
	}

	db.Callback().Create().Before("gorm:create").Register("otel:before_create", beforeFn("create"))
	db.Callback().Create().After("gorm:create").Register("otel:after_create", endSpan)

	db.Callback().Query().Before("gorm:query").Register("otel:before_query", beforeFn("query"))
	db.Callback().Query().After("gorm:query").Register("otel:after_query", endSpan)

	db.Callback().Update().Before("gorm:update").Register("otel:before_update", beforeFn("update"))
	db.Callback().Update().After("gorm:update").Register("otel:after_update", endSpan)

	db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete", beforeFn("delete"))
	db.Callback().Delete().After("gorm:delete").Register("otel:after_delete", endSpan)

	db.Callback().Row().Before("gorm:row").Register("otel:before_row", beforeFn("row"))
	db.Callback().Row().After("gorm:row").Register("otel:after_row", endSpan)

	db.Callback().Raw().Before("gorm:raw").Register("otel:before_raw", beforeFn("raw"))
	db.Callback().Raw().After("gorm:raw").Register("otel:after_raw", endSpan)
}

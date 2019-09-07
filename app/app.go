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

var (
	Logger *log.Logger
	DB     *gorm.DB
)

func Init() {
	Logger = getLogger()
	DB = getGormConnect()
}

func getLogger() *log.Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// ログ出力先を設定
	logPath := os.Getenv("LOG_PATH")
	logfile, err := os.OpenFile(logPath+"/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

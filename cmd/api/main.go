package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/yyh-gl/hobigon-golang-api-server/handler"
)

func main() {
	// ロガー設定
	logger := log.New(os.Stderr, "", log.LstdFlags)

	// ログ出力先を設定
	logfile, err := os.OpenFile("./logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open ./logs/app.log:" + err.Error())
	}
	defer logfile.Close()

	env := os.Getenv("APP_ENV")
	switch env {
	case "production":
		logger.SetOutput(logfile)
	default:
		logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	}

	fmt.Println("========================")
	fmt.Println("Server Start >> http://localhost:3000")
	fmt.Println("========================")
	logger.Fatal(http.ListenAndServe(":3000", handler.Router))
}

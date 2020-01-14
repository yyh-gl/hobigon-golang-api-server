package app

import (
	"io"
	"log"
	"os"
)

const (
	// APILogFilename : APIサーバ関連のログファイル名
	APILogFilename string = "api.log"
	// CLILogFilename : CLI関連のログファイル名
	CLILogFilename string = "cli.log"
)

// NewAPILogger : APIサーバ用のロガーを生成
func NewAPILogger() *log.Logger {
	logger := newLogger(APILogFilename)
	return logger
}

// NewCLILogger : CLI用のロガーを生成
func NewCLILogger() *log.Logger {
	logger := newLogger(CLILogFilename)
	return logger
}

// newLogger : 指定ファイルを出力先とするロガーを生成
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

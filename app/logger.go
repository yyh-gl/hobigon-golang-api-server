package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func NewLogger() {
	encoding := "json"
	if IsDev() || IsTest() {
		encoding = "console"
	}

	atomicLevel := zap.NewAtomicLevel()
	level := retrieveLeastLogLevel()
	atomicLevel.SetLevel(level)

	var err error
	logger, err = zap.Config{
		Level:       atomicLevel,
		Development: !IsPrd(),
		Encoding:    encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"/dev/stdout"},
		ErrorOutputPaths: []string{"/dev/stderr"},
	}.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func retrieveLeastLogLevel() zapcore.Level {
	if IsDev() || IsTest() {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func Error(err error) {
	logger.Error(err.Error(), zap.Error(err))
}

package logger

import "go.uber.org/zap"

type AppLogger struct {
	Log *zap.Logger
}

var Logger *AppLogger

func init() {
	Logger = &AppLogger{}
}

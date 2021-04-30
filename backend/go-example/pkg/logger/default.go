package logger

import "go.uber.org/zap"

func NewLogger() Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Failed to create a zap logger")
	}

	return logger.Sugar()
}

var Default = NewLogger()

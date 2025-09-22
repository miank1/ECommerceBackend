package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

// Init initializes the global logger
func Init() {
	var err error
	log, err = zap.NewProduction() // or zap.NewDevelopment() for local dev
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

// Get returns the global logger instance
func Get() *zap.Logger {
	if log == nil {
		Init()
	}
	return log
}

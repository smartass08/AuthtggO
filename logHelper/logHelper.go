package logHelper

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var gLogger *zap.SugaredLogger

func InitLogHelper(log_file string) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	core1 := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zap.InfoLevel)
	core2 := getFileLogger(log_file, cfg)
	logger := zap.New(zapcore.NewTee(core1, core2))
	defer logger.Sync() // flushes buffer, if any
	gLogger = logger.Sugar()
	gLogger.Info("Initialized LogHelper")
}

func getFileLogger(log_file string, cfg zapcore.EncoderConfig) zapcore.Core {
	//truncate old session's log file
	os.Remove(log_file)
	handleSync, _, err := zap.Open(log_file)
	if err != nil {
		log.Fatal(err)
	}
	return zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), handleSync, zap.InfoLevel)
}


func GetLogger() *zap.SugaredLogger {
	return gLogger
}
package clogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

// CreateLogger cria um Logger customizado.
func CreateLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	ts := strconv.FormatInt(time.Now().Unix(), 10)

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableStacktrace: false,
		DisableCaller:     true,
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
		OutputPaths: []string{
			"stdout",
			"logs/runtime_" + ts + ".log",
		},
		ErrorOutputPaths: []string{
			"stderr",
			"logs/error_" + ts + ".log",
		},
	}

	return zap.Must(config.Build())
}

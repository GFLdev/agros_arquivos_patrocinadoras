package main

import (
	"github.com/labstack/echo/v4"
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

// LogHTTPDetails registra os detalhes de uma requisição HTTP, incluindo o IP
// do cliente, o método e o caminho, com campos personalizados adicionais.
func LogHTTPDetails(
	c echo.Context,
	level zapcore.Level,
	msg string,
	fields ...zap.Field,
) {
	appCtx := GetAppContext(c)

	// Informações da requisição/resposta
	baseFields := []zap.Field{
		zap.String("method", c.Request().Method),
		zap.String("path", c.Path()),
		zap.String("client_ip", c.RealIP()),
	}
	allFields := append(baseFields, fields...)

	// Logger com o nível escolhido
	switch level {
	case zapcore.DebugLevel:
		appCtx.Logger.Debug(msg, allFields...)
	case zapcore.InfoLevel:
		appCtx.Logger.Info(msg, allFields...)
	case zapcore.WarnLevel:
		appCtx.Logger.Warn(msg, allFields...)
	case zapcore.ErrorLevel:
		appCtx.Logger.Error(msg, allFields...)
	case zapcore.FatalLevel:
		appCtx.Logger.Fatal(msg, allFields...)
	default:
		appCtx.Logger.Info(msg, allFields...)
	}
}

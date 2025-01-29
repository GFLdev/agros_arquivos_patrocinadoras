// Package logger fornece funcionalidades para criação e configuração de loggers
// utilizando a biblioteca zap. Ele permite a geração de logs formatados com
// diferentes níveis de severidade, personalização de encoders, além de suportar
// saída para console e arquivos.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

// CreateLogger cria e configura um logger utilizando a biblioteca zap.
//
// Retornos:
//   - *zap.Logger: instância configurada do logger.
func CreateLogger() *zap.Logger {
	// Configuração do encoder para produção
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Gerando um timestamp único para o nome do arquivo de log
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	// Configuração do logger
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
	}
	return zap.Must(config.Build())
}

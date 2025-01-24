package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

// CreateLogger cria e configura um zap.Logger para registrar informações de log
// em tempo de execução. O logger é configurado com um formato de codificação
// de console e um arquivo de log com timestamp no nome.
//
// Retorna o logger configurado.
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

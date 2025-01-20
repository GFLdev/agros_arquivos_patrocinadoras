package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

// BodyUnmarshall lê e desagrupa o body de uma requisição.
func BodyUnmarshall[T any](c echo.Context) (*T, error) {
	bodyJSON, err := io.ReadAll(c.Request().Body)
	if err != nil {
		LogHTTPDetails(c,
			zapcore.ErrorLevel,
			"Erro ao ler o corpo da requisição",
			zap.Error(err),
		)
		return nil, err
	}

	result := new(T)
	err = json.Unmarshal(bodyJSON, result)
	if err != nil {
		LogHTTPDetails(c,
			zapcore.ErrorLevel,
			"Erro ao desagrupar o corpo da requisição",
			zap.Error(err),
		)

		return nil, err
	}

	return result, nil
}

// GetAppContext recupera o AppContext do contexto do Echo.
func GetAppContext(c echo.Context) *AppContext {
	return c.Get("appContext").(*AppContext)
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

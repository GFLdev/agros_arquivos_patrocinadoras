package utils

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
)

// BodyUnmarshall lê e desagrupa o body de uma requisição.
func BodyUnmarshall[T any](c echo.Context, logr *zap.Logger) (*T, error) {
	bodyJSON, err := io.ReadAll(c.Request().Body)
	if err != nil {
		LogHTTPDetails(
			c,
			logr,
			zapcore.ErrorLevel,
			"Erro ao ler o corpo da requisição",
			zap.Error(err),
		)
		return nil, err
	}

	result := new(T)
	err = json.Unmarshal(bodyJSON, result)
	if err != nil {
		LogHTTPDetails(
			c,
			logr,
			zapcore.ErrorLevel,
			"Erro ao desagrupar o corpo da requisição",
			zap.Error(err),
		)

		return nil, err
	}

	return result, nil
}

// LogHTTPDetails registra os detalhes de uma requisição HTTP, incluindo o IP
// do cliente, o método e o caminho, com campos personalizados adicionais.
func LogHTTPDetails(
	c echo.Context,
	logr *zap.Logger,
	level zapcore.Level,
	msg string,
	fields ...zap.Field,
) {
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
		logr.Debug(msg, allFields...)
	case zapcore.InfoLevel:
		logr.Info(msg, allFields...)
	case zapcore.WarnLevel:
		logr.Warn(msg, allFields...)
	case zapcore.ErrorLevel:
		logr.Error(msg, allFields...)
	case zapcore.FatalLevel:
		logr.Fatal(msg, allFields...)
	default:
		logr.Info(msg, allFields...)
	}
}

func CheckAuthentication(c echo.Context) error {
	if false {
		c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		return c.JSON(http.StatusForbidden, ErrorRes{
			Message: "Não autorizado",
			Error:   fmt.Errorf("requisição por usuário não autorizado"),
		})
	}
	return nil
}

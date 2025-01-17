package main

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

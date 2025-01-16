package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
)

// BodyUnmarshall lê e desagrupa o body de uma requisição.
func BodyUnmarshall(c echo.Context, v interface{}) error {
	bodyJSON, err := io.ReadAll(c.Request().Body)
	if err != nil {
		LogMessage(
			c,
			"Erro ao ler o corpo da requisição de "+c.RealIP(),
			LogDetails{"Err": err.Error()},
		)

		return err
	}

	err = json.Unmarshal(bodyJSON, v)
	if err != nil {
		LogMessage(
			c,
			"Erro ao desagrupar o corpo da requisição de "+c.RealIP(),
			LogDetails{"Req": "\n" + string(bodyJSON), "Err": err.Error()},
		)

		return err
	}

	LogMessage(
		c,
		"Requisição recebida de "+c.RealIP(),
		LogDetails{"Req": "\n" + string(bodyJSON)},
	)

	return nil
}

package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	// Ler o corpo da requisição
	body := new(LoginRequest)
	err := BodyUnmarshall(c, body)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			ErrorResponse{
				Error:   "invalid request body",
				Message: err.Error(),
			},
		)
	}

	// TODO: Criar lógica de login
	res := LoginResponse{
		User:          "admin",
		Authenticated: true,
	}

	c.Response().Header().Add("Content-Type", "application/json")
	err = c.JSON(http.StatusOK, res)
	if err != nil {
		LogMessage(
			c,
			"Erro ao enviar resposta para "+c.RealIP(),
			LogDetails{"Res": body, "Err": err.Error()},
		)

		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{
				Error:   "failed to send response",
				Message: err.Error(),
			})
	}

	LogMessage(
		c,
		"Resposta enviada a "+c.RealIP(),
		LogDetails{"Res": body},
	)

	return nil
}

// DownloadFileHandler gerencia o processo de download de um arquivo.
func DownloadFileHandler(c echo.Context) error {
	// Ler o corpo da requisição
	body := new(DownloadRequest)
	err := BodyUnmarshall(c, body)
	if err != nil {
		return err
	}
	// TODO: Criar lógica de download de arquivo

	return nil
}

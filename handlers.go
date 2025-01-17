package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	// Ler o corpo da requisição
	body, err := BodyUnmarshall[LoginRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorResponse{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Validação
	if len(body.Username) == 0 || len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest,
			ErrorResponse{
				Message: "Usuário e/ou senha estão vazios",
				Error:   err,
			},
		)
	}

	// TODO: Criar lógica de login
	res := LoginResponse{
		User:          body.Username,
		Authenticated: true,
	}

	c.Response().Header().Add("Content-Type", "application/json")
	err = c.JSON(http.StatusOK, res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorResponse{
				Message: "Falha na resposta",
				Error:   err,
			})
	}

	return nil
}

package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	// Ler o corpo da requisição
	body, err := BodyUnmarshall[LoginReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// TODO: Criar lógica de login
	res := LoginRes{
		User:          body.Username,
		Authenticated: true,
	}

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	// Ler o corpo da requisição
	body, err := utils.BodyUnmarshall[utils.LoginReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			utils.ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// TODO: Criar lógica de login
	res := utils.LoginRes{
		User:          body.Username,
		Authenticated: true,
	}

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

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

	// Validação
	if len(body.Username) < 4 || len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Usuário e/ou senha estão vazios",
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

// AllUsersHandler obtém todas as categorias de arquivos de um
// usuário específico.
func AllUsersHandler(c echo.Context) error {
	ctx := GetAppContext(c)
	// TODO: Verificar se usuário é autenticado

	ctx.Repo.Lock()
	res := ctx.Repo.GetAllUsers()
	ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

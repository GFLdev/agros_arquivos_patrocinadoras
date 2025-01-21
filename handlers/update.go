package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UpdateUserHandler gerencia a modificação de um usuário pelo seu ID.
func UpdateUserHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[NameInputReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Atualização
	ctx.Repo.Lock()
	ok := ctx.Repo.UpdateUserById(userId, body.Name)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   fmt.Errorf("id de usuário não encontrado"),
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Usuário alterado com sucesso",
	})
}

// UpdateCategoryHandler gerencia a modificação de uma categoria pelo seu ID.
func UpdateCategoryHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[NameInputReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Atualização
	ctx.Repo.Lock()
	ok := ctx.Repo.UpdateCategoryById(userId, categId, body.Name)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não encontrada",
			Error:   fmt.Errorf("id de categoria não encontrado"),
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Categoria alterada com sucesso",
	})
}

// UpdateFileHandler gerencia a modificação de um arquivo pelo seu ID.
func UpdateFileHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[NameInputReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Atualização
	ctx.Repo.Lock()
	ok := ctx.Repo.UpdateFileById(userId, categId, fileId, body.Name)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   fmt.Errorf("id de usuário não encontrado"),
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Arquivo alterado com sucesso",
	})
}

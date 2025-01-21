package handlers

import (
	"github.com/google/uuid"
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteUserHandler gerencia a exclusão de um usuário pelo seu ID.
func DeleteUserHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	// Exclusão
	ctx.Repo.Lock()
	err := ctx.Repo.DeleteUserById(userId)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Usuário removido com sucesso",
	})
}

// DeleteCategoryHandler gerencia a exclusão de uma categoria pelo seu ID.
func DeleteCategoryHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	// Exclusão
	ctx.Repo.Lock()
	err := ctx.Repo.DeleteCategoryById(userId, categId)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não encontrada",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Categoria removida com sucesso",
	})
}

// DeleteFileHandler gerencia a exclusão de um arquivo pelo seu ID.
func DeleteFileHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	// Exclusão
	ctx.Repo.Lock()
	err := ctx.Repo.DeleteFileById(userId, categId, fileId)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Arquivo não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Arquivo removido com sucesso",
	})
}

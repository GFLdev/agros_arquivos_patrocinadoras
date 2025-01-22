package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"github.com/google/uuid"
	"net/http"

	"github.com/labstack/echo/v4"
)

// DeleteUserHandler gerencia a exclusão de um usuário pelo seu ID.
func DeleteUserHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	// Exclusão
	ctx.FSServ.Mux.Lock()
	err := ctx.FSServ.FS.DeleteUserById(userId)
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Usuário removido com sucesso",
	})
}

// DeleteCategoryHandler gerencia a exclusão de uma categoria pelo seu ID.
func DeleteCategoryHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	// Exclusão
	ctx.FSServ.Mux.Lock()
	err := ctx.FSServ.FS.DeleteCategoryById(userId, categId)
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Categoria não encontrada",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Categoria removida com sucesso",
	})
}

// DeleteFileHandler gerencia a exclusão de um arquivo pelo seu ID.
func DeleteFileHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	// Exclusão
	ctx.FSServ.Mux.Lock()
	err := ctx.FSServ.FS.DeleteFileById(userId, categId, fileId)
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Arquivo não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Arquivo removido com sucesso",
	})
}

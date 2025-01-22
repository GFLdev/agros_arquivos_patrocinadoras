package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// AllUsersHandler obtém todos os usuários do repositório.
func AllUsersHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	ctx.FSServ.Mux.Lock()
	res, ok := ctx.FSServ.FS.GetAllUsers()
	if !ok {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Nenhum usuário encontrado",
			Error:   fmt.Errorf("repositório não tem nenhum usuário"),
		})
	}
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// UserByIdHandler obtém um usuário com base em seu ID.
func UserByIdHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	ctx.FSServ.Mux.Lock()
	res, ok := ctx.FSServ.FS.GetUserById(userId)
	if !ok {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Usuário não encontrado",
			Error:   fmt.Errorf("usuário %s não encontrado", userId.String()),
		})
	}
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// AllCategoriesHandler obtém todas as categorias de um usuário do repositório.
func AllCategoriesHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	ctx.FSServ.Mux.Lock()
	res, ok := ctx.FSServ.FS.GetAllCategories(userId)
	if !ok {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Nenhuma categoria encontrada",
			Error:   fmt.Errorf("usuário %s não tem nenhuma categoria", userId.String()),
		})
	}
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// CategoryByIdHandler obtém uma categoria com base em seu ID.
func CategoryByIdHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	ctx.FSServ.Mux.Lock()
	res, ok := ctx.FSServ.FS.GetCategoryById(userId, categId)
	if !ok {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Categoria não encontrada",
			Error: fmt.Errorf("usuário %s não tem categoria %s",
				userId.String(),
				categId.String(),
			),
		})
	}
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// AllFilesHandler obtém todos os arquivos de uma categoria de um usuário do
// repositório.
func AllFilesHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	ctx.FSServ.Mux.Lock()
	res, ok := ctx.FSServ.FS.GetAllFiles(userId, categId)
	if !ok {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Nenhum arquivo encontrado",
			Error: fmt.Errorf("categoria %s de usuário %s não tem nenhum"+
				" arquivo",
				userId.String(),
				categId.String(),
			),
		})
	}
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// FileByIdHandler obtém um arquivo com base em seu ID.
func FileByIdHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	ctx.FSServ.Mux.Lock()
	res, ok := ctx.FSServ.FS.GetFileById(userId, categId, fileId)
	if !ok {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Arquivo não encontrado",
			Error: fmt.Errorf("categoria %s de usuário %s não tem arquivo %s",
				userId.String(),
				categId.String(),
				fileId.String(),
			),
		})
	}
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

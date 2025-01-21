package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// AllUsersHandler obtém todos os usuários do repositório.
func AllUsersHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	ctx.Repo.Lock()
	res, ok := ctx.Repo.GetAllUsers()
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum usuário encontrado",
			Error:   fmt.Errorf("repositório não tem nenhum usuário"),
		})
	}
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

// UserByIdHandler obtém um usuário com base em seu ID.
func UserByIdHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	ctx.Repo.Lock()
	res, ok := ctx.Repo.GetUserById(userId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   fmt.Errorf("usuário %s não encontrado", userId.String()),
		})
	}
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

// AllCategoriesHandler obtém todas as categorias de um usuário do repositório.
func AllCategoriesHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	ctx.Repo.Lock()
	res, ok := ctx.Repo.GetAllCategories(userId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhuma categoria encontrada",
			Error:   fmt.Errorf("usuário %s não tem nenhuma categoria", userId.String()),
		})
	}
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

// CategoryByIdHandler obtém uma categoria com base em seu ID.
func CategoryByIdHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	ctx.Repo.Lock()
	res, ok := ctx.Repo.GetCategoryById(userId, categId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não encontrada",
			Error: fmt.Errorf("usuário %s não tem categoria %s",
				userId.String(),
				categId.String(),
			),
		})
	}
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

// AllFilesHandler obtém todos os arquivos de uma categoria de um usuário do
// repositório.
func AllFilesHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	ctx.Repo.Lock()
	res, ok := ctx.Repo.GetAllFiles(userId, categId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum arquivo encontrado",
			Error: fmt.Errorf("categoria %s de usuário %s não tem nenhum"+
				" arquivo",
				userId.String(),
				categId.String(),
			),
		})
	}
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

// FileByIdHandler obtém um arquivo com base em seu ID.
func FileByIdHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	ctx.Repo.Lock()
	res, ok := ctx.Repo.GetFileById(userId, categId, fileId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Arquivo não encontrado",
			Error: fmt.Errorf("categoria %s de usuário %s não tem arquivo %s",
				userId.String(),
				categId.String(),
				fileId.String(),
			),
		})
	}
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

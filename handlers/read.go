package handlers

import (
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
	res := ctx.Repo.GetAllUsers()
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
	res := ctx.Repo.GetUserById(userId)
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
	res := ctx.Repo.GetAllCategories(userId)
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
	res := ctx.Repo.GetCategoryById(userId, categId)
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
	res := ctx.Repo.GetAllFiles(userId, categId)
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
	res := ctx.Repo.GetFileById(userId, categId, fileId)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	return c.JSON(http.StatusOK, res)
}

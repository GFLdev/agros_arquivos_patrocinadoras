package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateUserHandler gerencia a criação de um usuário.
func CreateUserHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

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

	ctx.Repo.Lock()
	err = ctx.Repo.CreateUser(body.Name)
	defer ctx.Repo.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar usuário",
				Error:   err,
			},
		)
	}

	return c.JSON(http.StatusOK, GenericRes{
		"Usuário criado com sucesso",
	})
}

// CreateCategoryHandler gerencia a criação de uma categoria.
func CreateCategoryHandler(c echo.Context) error {
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

	ctx.Repo.Lock()
	err = ctx.Repo.CreateCategory(userId, body.Name)
	defer ctx.Repo.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar categoria",
				Error:   err,
			},
		)
	}

	return c.JSON(http.StatusOK, GenericRes{
		"Categoria criada com sucesso",
	})
}

// CreateFileHandler gerencia a criação de um arquivo.
func CreateFileHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[FileInputReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	ctx.Repo.Lock()
	err = ctx.Repo.CreateFile(
		userId,
		categId,
		body.Name,
		body.FileType,
		body.Content,
	)
	defer ctx.Repo.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar arquivo",
				Error:   err,
			},
		)
	}

	return c.JSON(http.StatusOK, GenericRes{
		"Arquivo criado com sucesso",
	})
}

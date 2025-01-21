package handlers

import (
	"agros_arquivos_patrocinadoras/db"
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

	params := db.UpdateUserParams{
		UserId: uuid.MustParse(c.Param("userId")),
		Name:   body.Name,
	}

	// Atualização
	ctx.Repo.Lock()
	err = ctx.Repo.UpdateUserById(params)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
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

	params := db.UpdateCategoryParams{
		UserId:  uuid.MustParse(c.Param("userId")),
		CategId: uuid.MustParse(c.Param("categId")),
		Name:    body.Name,
	}

	// Atualização
	ctx.Repo.Lock()
	err = ctx.Repo.UpdateCategoryById(params)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não encontrada",
			Error:   err,
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

	params := db.UpdateFileParams{
		UserId:   uuid.MustParse(c.Param("userId")),
		CategId:  uuid.MustParse(c.Param("categId")),
		FileId:   uuid.MustParse(c.Param("fileId")),
		Name:     body.Name,
		FileType: body.FileType,
		Content:  body.Content,
	}

	// Atualização
	ctx.Repo.Lock()
	err = ctx.Repo.UpdateFileById(params)
	defer ctx.Repo.Unlock()

	c.Response().Header().Add("Content-Type", "application/json")
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Arquivo alterado com sucesso",
	})
}

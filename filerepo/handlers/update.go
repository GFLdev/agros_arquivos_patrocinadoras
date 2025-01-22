package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// UpdateUserHandler gerencia a modificação de um usuário pelo seu ID.
func UpdateUserHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	// Ler o corpo da requisição
	body, err := utils.BodyUnmarshall[utils.NameInputReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			utils.ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	params := fs.UpdateUserParams{
		UserId: uuid.MustParse(c.Param("userId")),
		Name:   body.Name,
	}

	// Atualização
	ctx.FSServ.Mux.Lock()
	err = ctx.FSServ.FS.UpdateUserById(params)
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Usuário alterado com sucesso",
	})
}

// UpdateCategoryHandler gerencia a modificação de uma categoria pelo seu ID.
func UpdateCategoryHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	// Ler o corpo da requisição
	body, err := utils.BodyUnmarshall[utils.NameInputReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			utils.ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	params := fs.UpdateCategoryParams{
		UserId:  uuid.MustParse(c.Param("userId")),
		CategId: uuid.MustParse(c.Param("categId")),
		Name:    body.Name,
	}

	// Atualização
	ctx.FSServ.Mux.Lock()
	err = ctx.FSServ.FS.UpdateCategoryById(params)
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Categoria não encontrada",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Categoria alterada com sucesso",
	})
}

// UpdateFileHandler gerencia a modificação de um arquivo pelo seu ID.
func UpdateFileHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	// Ler o corpo da requisição
	body, err := utils.BodyUnmarshall[utils.FileInputReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			utils.ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	params := fs.UpdateFileParams{
		UserId:    uuid.MustParse(c.Param("userId")),
		CategId:   uuid.MustParse(c.Param("categId")),
		FileId:    uuid.MustParse(c.Param("fileId")),
		Name:      body.Name,
		FileType:  body.FileType,
		Extension: body.Extension,
		Content:   body.Content,
	}

	// Atualização
	ctx.FSServ.Mux.Lock()
	err = ctx.FSServ.FS.UpdateFileById(params)
	defer ctx.FSServ.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Arquivo alterado com sucesso",
	})
}

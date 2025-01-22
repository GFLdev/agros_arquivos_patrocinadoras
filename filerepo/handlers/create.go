package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateUserHandler gerencia a criação de um usuário.
func CreateUserHandler(c echo.Context) error {
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

	req := fs.CreateUserParams{Name: body.Name}

	ctx.FSServ.Mux.Lock()
	err = ctx.FSServ.FS.CreateUser(req)
	defer ctx.FSServ.Mux.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.ErrorRes{
				Message: "Erro ao criar usuário",
				Error:   err,
			},
		)
	}

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, utils.GenericRes{
		Message: "Usuário criado com sucesso",
	})
}

// CreateCategoryHandler gerencia a criação de uma categoria.
func CreateCategoryHandler(c echo.Context) error {
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

	req := fs.CreateCategoryParams{
		UserId: uuid.MustParse(c.Param("userId")),
		Name:   body.Name,
	}

	ctx.FSServ.Mux.Lock()
	err = ctx.FSServ.FS.CreateCategory(req)
	defer ctx.FSServ.Mux.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.ErrorRes{
				Message: "Erro ao criar categoria",
				Error:   err,
			},
		)
	}

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, utils.GenericRes{
		"Categoria criada com sucesso",
	})
}

// CreateFileHandler gerencia a criação de um arquivo.
func CreateFileHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

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

	req := fs.CreateFileParams{
		UserId:    userId,
		CategId:   categId,
		Name:      body.Name,
		FileType:  body.FileType,
		Extension: body.Extension,
		Content:   body.Content,
	}

	ctx.FSServ.Mux.Lock()
	err = ctx.FSServ.FS.CreateFile(req)
	defer ctx.FSServ.Mux.Unlock()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.ErrorRes{
				Message: "Erro ao criar arquivo",
				Error:   err,
			},
		)
	}

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, utils.GenericRes{
		"Arquivo criado com sucesso",
	})
}

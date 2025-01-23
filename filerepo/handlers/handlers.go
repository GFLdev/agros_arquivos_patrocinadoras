package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/auth"
	"agros_arquivos_patrocinadoras/filerepo/db"
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := services.GetContext(c)

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[LoginReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// TODO: Criar lógica de login
	res := LoginRes{
		User:          body.Username,
		Authenticated: true,
	}

	return c.JSON(http.StatusOK, res)
}

// DownloadHandler gerencia o envio de um arquivo para download.
func DownloadHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := services.GetContext(c)

	// Autenticação
	if auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Errorf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err,
			},
		)
	}

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err,
			},
		)
	}

	// Obtenção dos metadados do arquivo
	file, err := db.QueryFileById(ctx.DB, fileId)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Arquivo não encontrado",
				Error:   err,
			},
		)
	}

	// Anexar arquivo
	attach, err := ctx.FileSystem.GetFile(
		userId,
		categId,
		fileId,
		file.Extension,
	)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Arquivo não encontrado",
				Error:   err,
			},
		)
	}

	// Cabeçalho para o arquivo e resposta
	c.Response().Header().Add(echo.HeaderContentType, file.Mimetype)
	return c.Attachment(attach, file.Name)
}

// ----------
//   CREATE
// ----------

// CreateUserHandler gerencia a criação de um usuário.
func CreateUserHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := services.GetContext(c)

	// Autenticação
	if auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Errorf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[CreateUserReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Geração do Id
	userId := uuid.New()

	// Criação do diretório
	err = ctx.FileSystem.CreateUser(userId)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar diretório do usuário",
				Error:   err,
			},
		)
	}

	// Agendar a exclusão do diretório em caso de erro subsequente
	defer func() {
		if err != nil {
			err = ctx.FileSystem.DeleteUser(userId)
			if err != nil {
				ctx.Logger.Error(
					"Erro ao limpar arquivo lixo",
					zap.Error(err),
				)
			}
		}
	}()

	// Criptografia
	hash, err := HashPassword(body.Password)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criptografar senha",
				Error:   err,
			},
		)
	}

	// Criar usuário no banco
	err = db.CreateUser(ctx.DB, userId, body.Name, hash)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar usuário",
				Error:   err,
			},
		)
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Usuário criado com sucesso",
	})
}

// CreateCategoryHandler gerencia a criação de uma categoria.
func CreateCategoryHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := services.GetContext(c)

	// Autenticação
	if auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Errorf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[CreateCategoryReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	// Geração do Id
	categId := uuid.New()

	// Criação do diretório
	err = ctx.FileSystem.CreateCategory(userId, categId)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar diretório da categoria",
				Error:   err,
			},
		)
	}

	// Agendar a exclusão do diretório em caso de erro subsequente
	defer func() {
		if err != nil {
			err = ctx.FileSystem.DeleteCategory(userId, categId)
			if err != nil {
				ctx.Logger.Error(
					"Erro ao limpar arquivo lixo",
					zap.Error(err),
				)
			}
		}
	}()

	// Criar categoria no banco
	err = db.CreateCategory(ctx.DB, categId, userId, body.Name)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
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
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := services.GetContext(c)

	// Autenticação
	if auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Errorf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[CreateFileReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err,
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err,
			},
		)
	}

	// Geração do Id
	fileId := uuid.New()

	// Criação do diretório
	err = ctx.FileSystem.CreateFile(
		userId,
		categId,
		fileId,
		body.Extension,
		&body.Content,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar arquivo",
				Error:   err,
			},
		)
	}

	// Agendar a exclusão do diretório em caso de erro subsequente
	defer func() {
		if err != nil {
			err = ctx.FileSystem.DeleteFile(
				userId,
				categId,
				fileId,
				body.Extension,
			)
			if err != nil {
				ctx.Logger.Error(
					"Erro ao limpar arquivo lixo",
					zap.Error(err),
				)
			}
		}
	}()

	// Criar categoria no banco
	err = db.CreateFile(
		ctx.DB,
		fileId,
		categId,
		body.Name,
		body.Extension,
		body.Mimetype,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
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

// --------
//   READ
// --------

// AllUsersHandler obtém todos os usuários do repositório.
func AllUsersHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	ctx.FileSystem.Mux.Lock()
	res, ok := ctx.FileSystem.FS.GetAllUsers()
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum usuário encontrado",
			Error:   fmt.Errorf("repositório não tem nenhum usuário"),
		})
	}
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// UserByIdHandler obtém um usuário com base em seu ID.
func UserByIdHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	ctx.FileSystem.Mux.Lock()
	res, ok := ctx.FileSystem.FS.GetUserById(userId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   fmt.Errorf("usuário %s não encontrado", userId.String()),
		})
	}
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// AllCategoriesHandler obtém todas as categorias de um usuário do repositório.
func AllCategoriesHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	ctx.FileSystem.Mux.Lock()
	res, ok := ctx.FileSystem.FS.GetAllCategories(userId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhuma categoria encontrada",
			Error:   fmt.Errorf("usuário %s não tem nenhuma categoria", userId.String()),
		})
	}
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// CategoryByIdHandler obtém uma categoria com base em seu ID.
func CategoryByIdHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	ctx.FileSystem.Mux.Lock()
	res, ok := ctx.FileSystem.FS.GetCategoryById(userId, categId)
	if !ok {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não encontrada",
			Error: fmt.Errorf("usuário %s não tem categoria %s",
				userId.String(),
				categId.String(),
			),
		})
	}
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// AllFilesHandler obtém todos os arquivos de uma categoria de um usuário do
// repositório.
func AllFilesHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	ctx.FileSystem.Mux.Lock()
	res, ok := ctx.FileSystem.FS.GetAllFiles(userId, categId)
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
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// FileByIdHandler obtém um arquivo com base em seu ID.
func FileByIdHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	ctx.FileSystem.Mux.Lock()
	res, ok := ctx.FileSystem.FS.GetFileById(userId, categId, fileId)
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
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(http.StatusOK, res)
}

// ----------
//   UPDATE
// ----------

// UpdateUserHandler gerencia a modificação de um usuário pelo seu ID.
func UpdateUserHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[utils.NameInputReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
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
	ctx.FileSystem.Mux.Lock()
	err = ctx.FileSystem.FS.UpdateUserById(params)
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[utils.NameInputReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
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
	ctx.FileSystem.Mux.Lock()
	err = ctx.FileSystem.FS.UpdateCategoryById(params)
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[FileInputReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
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
	ctx.FileSystem.Mux.Lock()
	err = ctx.FileSystem.FS.UpdateFileById(params)
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

// ----------
//   DELETE
// ----------

// DeleteUserHandler gerencia a exclusão de um usuário pelo seu ID.
func DeleteUserHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))

	// Exclusão
	ctx.FileSystem.Mux.Lock()
	err := ctx.FileSystem.FS.DeleteUserById(userId)
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))

	// Exclusão
	ctx.FileSystem.Mux.Lock()
	err := ctx.FileSystem.FS.DeleteCategoryById(userId, categId)
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	ctx := services.GetContext(c)

	if err := auth.IsAuthenticated(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	// Exclusão
	ctx.FileSystem.Mux.Lock()
	err := ctx.FileSystem.FS.DeleteFileById(userId, categId, fileId)
	defer ctx.FileSystem.Mux.Unlock()

	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
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

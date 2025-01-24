package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/auth"
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

	// Criar usuário
	err = db.CreateUser(ctx.DB, ctx.FileSystem, userId, body.Name, hash)
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

// GetAllUsers obtém todos os usuários do repositório.
func GetAllUsers(c echo.Context) error {
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

	// Obtenção dos dados
	res, err := db.QueryAllUsers(ctx.DB)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum usuário foi obtido",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserById obtém um usuário com base em seu Id.
func GetUserById(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	// Obtenção dos dados
	res, err := db.QueryUserById(ctx.DB, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não obtido",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetAllCategories obtém todas as categorias de um usuário do repositório.
func GetAllCategories(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	// Obtenção dos dados
	res, err := db.QueryAllCategories(ctx.DB, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhuma categoria foi obtida",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetCategoryById obtém uma categoria com base em seu Id.
func GetCategoryById(c echo.Context) error {
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
	_, err := uuid.Parse(c.Param("userId"))
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

	// Obtenção dos dados
	res, err := db.QueryCategoryById(ctx.DB, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não obtida",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetAllFiles obtém todos os arquivos de uma categoria de um usuário do
// repositório.
func GetAllFiles(c echo.Context) error {
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
	_, err := uuid.Parse(c.Param("userId"))
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

	// Obtenção dos dados
	res, err := db.QueryAllFiles(ctx.DB, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum arquivo foi obtido",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetFileById obtém um arquivo com base em seu Id.
func GetFileById(c echo.Context) error {
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
	_, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	_, err = uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err,
			},
		)
	}

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err,
			},
		)
	}

	// Obtenção dos dados
	res, err := db.QueryFileById(ctx.DB, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Arquivo não obtido",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// ----------
//   UPDATE
// ----------

// UpdateUserHandler gerencia a modificação de um usuário pelo seu Id.
func UpdateUserHandler(c echo.Context) error {
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
	body, err := BodyUnmarshall[UpdateUserReq](c, ctx.Logger)
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

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Usuário alterado com sucesso",
	})
}

// UpdateCategoryHandler gerencia a modificação de uma categoria pelo seu Id.
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

// UpdateFileHandler gerencia a modificação de um arquivo pelo seu Id.
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

// DeleteUser gerencia a exclusão de um usuário pelo seu Id.
func DeleteUser(c echo.Context) error {
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
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err,
			},
		)
	}

	// Remoção do diretório
	err = ctx.FileSystem.DeleteUser(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{
			Message: "Não foi possível excluir diretório do usuário",
			Error:   err,
		})
	}

	// Remoção do usuário
	err = db.DeleteUser(ctx.DB, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não encontrado",
			Error:   err,
		})
	}

	return c.JSON(
		http.StatusOK,
		GenericRes{Message: "Usuário removido com sucesso"},
	)
}

// DeleteCategory gerencia a exclusão de uma categoria pelo seu Id.
func DeleteCategory(c echo.Context) error {
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

	// Remoção da categoria
	err = db.DeleteCategory(ctx.DB, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não encontrada",
			Error:   err,
		})
	}

	// Remoção do diretório
	err = ctx.FileSystem.DeleteCategory(userId, categId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{
			Message: "Não foi possível excluir diretório da categoria",
			Error:   err,
		})
	}

	return c.JSON(
		http.StatusOK,
		GenericRes{Message: "Categoria removida com sucesso"},
	)
}

// DeleteFile gerencia a exclusão de um arquivo pelo seu Id.
func DeleteFile(c echo.Context) error {
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

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err,
			},
		)
	}

	// Obtenção dos dados
	file, err := db.QueryFileById(ctx.DB, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Arquivo não obtido",
			Error:   err,
		})
	}

	// Remoção do arquivo
	err = db.DeleteFile(ctx.DB, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Arquivo não encontrado",
			Error:   err,
		})
	}

	// Remoção do arquivo em disco
	err = ctx.FileSystem.DeleteFile(
		userId,
		categId,
		fileId,
		file.Extension,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{
			Message: "Não foi possível excluir diretório da categoria",
			Error:   err,
		})
	}

	return c.JSON(
		http.StatusOK,
		GenericRes{Message: "Arquivo removido com sucesso"},
	)
}

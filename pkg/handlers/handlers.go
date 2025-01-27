package handlers

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/store"
	"agros_arquivos_patrocinadoras/pkg/auth"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[LoginReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
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
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obtenção dos metadados do arquivo
	file, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Arquivo não encontrado",
				Error:   err.Error(),
			},
		)
	}

	// Anexar arquivo
	path := fmt.Sprintf(
		"%s/%s/%s/%s%s",
		ctx.FileSystem.Root,
		userId,
		categId,
		fileId,
		file.Extension)
	exists := ctx.FileSystem.EntityExists(path)
	if !exists {
		return c.JSON(
			http.StatusNotFound,
			ErrorRes{
				Message: "Arquivo não encontrado",
				Error:   fmt.Sprintf("arquivo em %s não encontrado", path),
			},
		)
	}

	// Cabeçalho para o arquivo e resposta
	c.Response().Header().Add(echo.HeaderContentType, file.Mimetype)
	return c.Attachment(path, file.Name)
}

// ----------
//   CREATE
// ----------

// CreateUserHandler gerencia a criação de um usuário.
func CreateUserHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
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
				Error:   err.Error(),
			},
		)
	}

	// Criptografia
	hash, err := HashPassword(body.Password)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criptografar senha",
				Error:   err.Error(),
			},
		)
	}

	// Criar usuário
	user := store.UserParams{
		Name:     body.Name,
		Password: hash,
	}
	err = store.CreateUser(ctx, &user)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar usuário",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[CreateCategoryReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err.Error(),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	// Criar categoria
	categ := store.CategParams{
		UserId: userId,
		Name:   body.Name,
	}
	err = store.CreateCategory(ctx, &categ)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar categoria",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[CreateFileReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err.Error(),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	// Criar categoria no banco
	file := store.FileParams{
		UserId:    userId,
		CategId:   categId,
		Name:      body.Name,
		Extension: body.Extension,
		Mimetype:  body.Mimetype,
		Content:   &body.Content,
	}
	err = store.CreateFile(ctx, &file)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criar arquivo",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Obtenção dos dados
	res, err := store.QueryAllUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum usuário foi obtido",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserById obtém um usuário com base em seu Id.
func GetUserById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obtenção dos dados
	res, err := store.QueryUserById(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Usuário não obtido",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetAllCategories obtém todas as categorias de um usuário do repositório.
func GetAllCategories(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obtenção dos dados
	res, err := store.QueryAllCategories(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhuma categoria foi obtida",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetCategoryById obtém uma categoria com base em seu Id.
func GetCategoryById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	_, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obtenção dos dados
	res, err := store.QueryCategoryById(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Categoria não obtida",
			Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	_, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obtenção dos dados
	res, err := store.QueryAllFiles(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Nenhum arquivo foi obtido",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

// GetFileById obtém um arquivo com base em seu Id.
func GetFileById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	_, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	_, err = uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obtenção dos dados
	res, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorRes{
			Message: "Arquivo não obtido",
			Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[UpdateUserReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err.Error(),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obter dados do usuário
	data, err := store.QueryUserById(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound,
			ErrorRes{
				Message: "Usuário não obtido",
				Error:   err.Error(),
			})
	}

	// Criptografia
	hash, err := HashPassword(body.Password)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao criptografar senha",
				Error:   err.Error(),
			},
		)
	}

	// Alteração
	user := store.UserUpdate{
		UserId:  userId,
		OldName: data.Name,
		UserParams: store.UserParams{
			Name:     body.Name,
			Password: hash,
		},
	}
	err = store.UpdateUser(ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao atualizar usuário",
				Error:   err.Error(),
			})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Usuário alterado com sucesso",
	})
}

// UpdateCategoryHandler gerencia a modificação de uma categoria pelo seu Id.
func UpdateCategoryHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[UpdateCategoryReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err.Error(),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obter dados do usuário
	data, err := store.QueryCategoryById(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound,
			ErrorRes{
				Message: "Usuário não obtido",
				Error:   err.Error(),
			})
	}

	// Alteração
	categ := store.CategUpdate{
		CategId:   categId,
		OldUserId: userId,
		OldName:   data.Name,
		CategParams: store.CategParams{
			UserId: body.UserId,
			Name:   body.Name,
		},
	}
	err = store.UpdateCategory(ctx, categ)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao atualizar usuário",
				Error:   err.Error(),
			})
	}

	return c.JSON(http.StatusOK, GenericRes{
		Message: "Categoria alterada com sucesso",
	})
}

// UpdateFileHandler gerencia a modificação de um arquivo pelo seu Id.
func UpdateFileHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Ler o corpo da requisição
	body, err := BodyUnmarshall[UpdateFileReq](c, ctx.Logger)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Body da requisição inválido",
				Error:   err.Error(),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obter dados do usuário
	data, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound,
			ErrorRes{
				Message: "Usuário não obtido",
				Error:   err.Error(),
			})
	}

	// Alteração
	file := store.FileUpdate{
		FileId:       fileId,
		OldCategId:   categId,
		OldName:      data.Name,
		OldExtension: data.Extension,
		OldMimetype:  data.Mimetype,
		FileParams: store.FileParams{
			UserId:    userId,
			CategId:   body.CategId,
			Name:      body.Name,
			Extension: body.Extension,
			Mimetype:  body.Mimetype,
			Content:   &body.Content,
		},
	}
	err = store.UpdateFile(ctx, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao atualizar arquivo",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	// Remoção do usuário
	user := store.UserDelete{UserId: userId}
	err = store.DeleteUser(ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao excluir usuário",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	// Remoção da categoria
	categ := store.CategDelete{
		UserId:  userId,
		CategId: categId,
	}
	err = store.DeleteCategory(ctx, categ)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao excluir categoria",
				Error:   err.Error(),
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
	ctx := context.GetContext(c)

	// Autenticação
	if !auth.IsAuthenticated(c) {
		return c.JSON(
			http.StatusUnauthorized,
			ErrorRes{
				Message: "Usuário não autorizado",
				Error:   fmt.Sprintf("usuário não tem permissões para esta operação"),
			},
		)
	}

	// Parâmetros
	userId, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de usuário inválido",
				Error:   err.Error(),
			},
		)
	}

	categId, err := uuid.Parse(c.Param("categId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de categoria inválido",
				Error:   err.Error(),
			},
		)
	}

	fileId, err := uuid.Parse(c.Param("fileId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			ErrorRes{
				Message: "Id de arquivo inválido",
				Error:   err.Error(),
			},
		)
	}

	// Obter extensão do arquivo
	data, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound,
			ErrorRes{
				Message: "Arquivo não obtido",
				Error:   err.Error(),
			})
	}

	// Remoção do arquivo
	file := store.FileDelete{
		CategDelete: store.CategDelete{
			UserId:  userId,
			CategId: categId,
		},
		FileId:    fileId,
		Extension: data.Extension,
	}
	err = store.DeleteFile(ctx, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrorRes{
				Message: "Erro ao excluir arquivo",
				Error:   err.Error(),
			})
	}

	return c.JSON(
		http.StatusOK,
		GenericRes{Message: "Arquivo removido com sucesso"},
	)
}

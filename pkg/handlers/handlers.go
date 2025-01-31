// Package handlers contém os manipuladores HTTP responsáveis por processar
// as requisições recebidas pela aplicação. Ele inclui funcionalidades
// como autenticação, gerenciamento de usuários, categorias e arquivos.
// Este pacote utiliza o framework Echo para lidar com as requisições HTTP
// e interage com outras camadas da aplicação para realizar as operações
// necessárias.
package handlers

import (
	"agros_arquivos_patrocinadoras/pkg/app"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/auth"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoginHandler gerencia o processo de login de um usuário.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     seja bem-sucedido.
func LoginHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[LoginReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Verificar credenciais
	user := app.UserParams{
		Name:     body.Name,
		Password: body.Password,
	}
	userId, err := app.GetCredentials(ctx, user)
	if err != nil || userId == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Gerar token e resposta
	token, err := auth.GenerateToken(c, userId, body.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	c.Response().Header().Add(echo.HeaderAuthorization, "Bearer "+token)
	return c.JSON(http.StatusOK, echo.Map{
		"token":   token,
		"message": LoginSuccessMessage,
	})
}

// DownloadHandler gerencia o envio de um arquivo para download.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     seja bem-sucedido.
func DownloadHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obtenção do contexto da aplicação e dos metadados do arquivo
	ctx := context.GetContext(c)
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, FileNotFoundMessage)
	}
	path := fmt.Sprintf(
		"%s/%s/%s/%s%s",
		ctx.FileSystem.Root,
		userId,
		categId,
		fileId,
		file.Extension,
	)

	// Verificar existência no sistema de arquivos
	exists := ctx.FileSystem.EntityExists(path)
	if !exists {
		return c.JSON(http.StatusNotFound, FileNotFoundMessage)
	}
	c.Response().Header().Add(echo.HeaderContentType, file.Mimetype)
	return c.Attachment(path, file.Name)
}

// CreateUserHandler gerencia a criação de um novo usuário no sistema.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a criação
//     do usuário seja bem-sucedida.
func CreateUserHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateUserReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Validar senha
	if len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest, InvalidPasswordMessage)
	}
	user := app.UserParams{
		Name:     body.Name,
		Password: body.Password,
	}
	id, err := app.CreateUser(ctx, user)
	if err != nil {
		if err.Error() == "nome de usuário já existente" {
			return c.JSON(http.StatusConflict, DuplicateUserMessage)
		}
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"id":      id.String(),
		"message": string(CreatedUserMessage),
	})
}

// CreateCategoryHandler gerencia a criação de uma nova categoria associada
// a um usuário.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a criação
//     da categoria seja bem-sucedida.
func CreateCategoryHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateCategoryReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Criar categoria
	categ := app.CategParams{
		UserId: userId,
		Name:   body.Name,
	}
	if _, err = app.CreateCategory(ctx, categ); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusCreated, CreatedCategoryMessage)
}

// CreateFileHandler gerencia a criação de um novo arquivo em uma categoria
// existente.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a criação
//     do arquivo seja bem-sucedida.
func CreateFileHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateFileReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Criar arquivo
	file := app.FileParams{
		UserId:    userId,
		CategId:   categId,
		Name:      body.Name,
		Extension: body.Extension,
		Mimetype:  body.Mimetype,
		Content:   &body.Content,
	}
	if _, err = app.CreateFile(ctx, file); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusCreated, CreatedFileMessage)
}

// GetAllUsers obtém todos os usuários presentes no repositório.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a listagem
//     dos usuários seja bem-sucedida.
func GetAllUsers(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e de todos os usuários
	ctx := context.GetContext(c)
	res, err := app.QueryAllUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, UsersNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetUserById obtém um usuário específico com base em seu identificador único.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     seja bem-sucedido.
func GetUserById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obtenção do contexto da aplicação e do usuário
	ctx := context.GetContext(c)
	res, err := app.QueryUserById(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetAllCategories obtém todas as categorias presentes no repositório.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a listagem
//     das categorias seja bem-sucedida.
func GetAllCategories(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obtenção do contexto da aplicação e de todas as categorias
	ctx := context.GetContext(c)
	res, err := app.QueryAllCategories(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, CategoriesNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetCategoryById obtém uma categoria específica com base em seu identificador
// único.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     seja bem-sucedido.
func GetCategoryById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	_, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obtenção do contexto da aplicação e da categoria
	ctx := context.GetContext(c)
	res, err := app.QueryCategoryById(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetAllFiles obtém todos os arquivos presentes no repositório.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a listagem
//     dos arquivos seja bem-sucedida.
func GetAllFiles(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	_, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obtenção do contexto da aplicação e de todos os arquivos
	ctx := context.GetContext(c)
	res, err := app.QueryAllFiles(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, FilesNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetFileById obtém um arquivo específico com base em seu identificador único.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     seja bem-sucedido.
func GetFileById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros
	_, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	_, err = ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obtenção do contexto da aplicação e do arquivo
	ctx := context.GetContext(c)
	res, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateUserHandler gerencia a atualização dos dados de um usuário existente.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a atualização
//     dos dados seja bem-sucedida.
func UpdateUserHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateUserReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obter os dados do usuário
	data, err := app.QueryUserById(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	// Validar senha e alteração
	if len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}
	user := app.UserUpdate{
		UserId:  userId,
		OldName: data.Name,
		UserParams: app.UserParams{
			Name:     body.Name,
			Password: body.Password,
		},
	}
	if err = app.UpdateUser(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, UpdatedUserMessage)
}

// UpdateCategoryHandler gerencia a atualização dos dados de uma categoria
// existente.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a atualização
//     dos dados seja bem-sucedida.
func UpdateCategoryHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateCategoryReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obter dados do usuário
	data, err := app.QueryCategoryById(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	// Alteração
	categ := app.CategUpdate{
		CategId:   categId,
		OldUserId: userId,
		OldName:   data.Name,
		CategParams: app.CategParams{
			UserId: body.UserId,
			Name:   body.Name,
		},
	}
	if err = app.UpdateCategory(ctx, categ); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, UpdatedCategoryMessage)
}

// UpdateFileHandler gerencia a atualização dos dados de um arquivo existente.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso a atualização
//     dos dados seja bem-sucedida.
func UpdateFileHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateFileReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obter dados do usuário
	data, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	// Alteração
	file := app.FileUpdate{
		FileId:       fileId,
		OldCategId:   categId,
		OldName:      data.Name,
		OldExtension: data.Extension,
		OldMimetype:  data.Mimetype,
		FileParams: app.FileParams{
			UserId:    userId,
			CategId:   body.CategId,
			Name:      body.Name,
			Extension: body.Extension,
			Mimetype:  body.Mimetype,
			Content:   &body.Content,
		},
	}
	if err = app.UpdateFile(ctx, file); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	return c.JSON(http.StatusOK, UpdatedFileMessage)
}

// DeleteUser gerencia a exclusão de um usuário existente no sistema.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     de exclusão seja bem-sucedido.
func DeleteUser(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obtenção do contexto da aplicação e remoção do usuário
	ctx := context.GetContext(c)
	user := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedUserMessage)
}

// DeleteCategory gerencia a exclusão de uma categoria existente no sistema.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     de exclusão seja bem-sucedido.
func DeleteCategory(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obtenção do contexto da aplicação e remoção da categoria
	ctx := context.GetContext(c)
	categ := app.CategDelete{
		UserDelete: app.UserDelete{UserId: userId},
		CategId:    categId,
	}
	if err = app.DeleteCategory(ctx, categ); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedCategoryMessage)
}

// DeleteFile gerencia a exclusão de um arquivo existente no sistema.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     de exclusão seja bem-sucedido.
func DeleteFile(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obtenção do contexto da aplicação e dos dados do arquivo
	ctx := context.GetContext(c)
	data, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, FileNotFoundMessage)
	}

	// Remoção do arquivo
	file := app.FileDelete{
		CategDelete: app.CategDelete{
			UserDelete: app.UserDelete{UserId: userId},
			CategId:    categId,
		},
		FileId:    fileId,
		Extension: data.Extension,
	}
	if err = app.DeleteFile(ctx, file); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedFileMessage)
}

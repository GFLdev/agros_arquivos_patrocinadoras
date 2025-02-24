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
	"agros_arquivos_patrocinadoras/pkg/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// SessionHandler retorna os dados do JWT token para a sessão de um usuário.
//
// Parâmetros:
//   - c: contexto Echo contendo as informações da requisição HTTP.
//
// Retorno:
//   - error: um erro HTTP apropriado em caso de falha ou nil caso o processo
//     seja bem-sucedido.
func SessionHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação
	ctx := context.GetContext(c)

	claims, err := auth.GetClaims(c)
	if err != nil {
		ctx.Logger.Error("Erro ao obter claims.", zap.Error(err))
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Checar se é administrador
	admin := auth.AuthenticateAdmin(c)

	return c.JSON(http.StatusOK, echo.Map{
		"id":    claims.Id,
		"name":  claims.Name,
		"admin": admin,
	})
}

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

	// Validar parâmetros
	if body.Username == "" {
		return c.JSON(http.StatusBadRequest, EmptyUsernameMessage)
	}
	if len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest, InvalidPasswordMessage)
	}

	// Verificar credenciais
	loginParams := app.LoginParams{
		Username: body.Username,
		Password: body.Password,
	}
	loginData, err := app.QueryLogin(ctx, loginParams)
	if err != nil || loginData.UserId == uuid.Nil {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Verificar se é administrador
	admin := auth.AuthenticateAdmin(c)

	// Gerar token
	duration := time.Duration(ctx.Config.JwtExpires) * time.Minute
	expiresAt := time.Now().Add(duration)
	claimsData := auth.ClaimsData{
		Id:   loginData.UserId,
		Name: loginData.Name,
	}
	token, err := auth.GenerateToken(c, claimsData, expiresAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	// Adicionar cookie e resposta
	jwtCookie := &http.Cookie{
		Name:    "jwt",
		Value:   token,
		Expires: expiresAt,
	}
	c.SetCookie(jwtCookie)
	res := LoginRes{
		Token:   token,
		Message: LoginSuccessMessage,
		Id:      loginData.UserId.String(),
		Name:    loginData.Name,
		Admin:   admin,
	}
	return c.JSON(http.StatusOK, res)
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateUserReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Validar parâmetros
	if body.Username == "" {
		return c.JSON(http.StatusBadRequest, EmptyUsernameMessage)
	}
	if body.Name == "" {
		return c.JSON(http.StatusBadRequest, EmptyNameMessage)
	}
	if len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest, InvalidPasswordMessage)
	}

	// Criar usuário
	user := app.UserData{
		Username: body.Username,
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

	// Resposta
	res := CreateResponse{
		Id:      id,
		Message: CreatedUserMessage,
	}
	return c.JSON(http.StatusCreated, res)
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateCategoryReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	// Criar categoria
	categ := app.CategData{
		UserId: userId,
		Name:   body.Name,
	}
	id, err := app.CreateCategory(ctx, categ)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	// Resposta
	res := CreateResponse{
		Id:      id,
		Message: CreatedCategoryMessage,
	}
	return c.JSON(http.StatusCreated, res)
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateFileReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	// Criar arquivo
	file := app.FileData{
		CategId:   categId,
		Name:      body.Name,
		Extension: body.Extension,
		Mimetype:  body.Mimetype,
		Content:   &body.Content,
	}
	id, err := app.CreateFile(ctx, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	// Resposta
	res := CreateResponse{
		Id:      id,
		Message: CreatedFileMessage,
	}
	return c.JSON(http.StatusCreated, res)
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Autorizar usuário
	if check := auth.AuthenticateUser(c, userId); !check {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do usuário
	user, err := app.QueryUserById(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}
	return c.JSON(http.StatusOK, user)
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	// Autorizar usuário
	if check := auth.AuthenticateUser(c, userId); !check {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção de todas as categorias
	categs, err := app.QueryAllCategories(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, CategoriesNotFoundMessage)
	}
	return c.JSON(http.StatusOK, categs)
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Autorizar usuário
	if check := auth.AuthenticateUser(c, userId); !check {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção da categoria
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}
	return c.JSON(http.StatusOK, categ)
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	// Autorizar usuário
	if check := auth.AuthenticateUser(c, userId); !check {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção de todos os arquivos
	files, err := app.QueryAllFiles(ctx, categId)
	if err != nil {
		return c.JSON(http.StatusNotFound, FilesNotFoundMessage)
	}
	return c.JSON(http.StatusOK, files)
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	fileId, err := ParseEntityUUID(c, File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Autorizar usuário
	if check := auth.AuthenticateUser(c, userId); !check {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção dos dados do arquivo
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		return c.JSON(http.StatusNotFound, FileNotFoundMessage)
	}
	return c.JSON(http.StatusOK, file)
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateUserReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Caso nada seja requisitado para alterar
	if body.Username == "" && body.Name == "" && body.Password == "" {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL e verificar se usuário existe
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	// Validar senha e alteração
	if body.Password != "" && len(body.Password) < 4 {
		return c.JSON(http.StatusBadRequest, InvalidPasswordMessage)
	}
	userParams := app.UserData{
		Username: body.Username,
		Name:     body.Name,
		Password: body.Password,
	}
	if err = app.UpdateUser(ctx, userId, userParams); err != nil {
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateCategoryReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Caso nada seja requisitado para alterar
	if body.UserId == "" && body.Name == "" {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Transformar UserId para atualização e verificar existência
	parsedUserId := uuid.Nil
	if body.UserId != "" {
		parsedUserId, err = uuid.Parse(body.UserId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
		}
		if _, err = app.QueryUserById(ctx, parsedUserId); err != nil {
			return c.JSON(http.StatusNotFound, UserNotFoundMessage)
		}
	}

	// Parâmetros da URL e verificar se usuário e categoria existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	// Alteração
	categParams := app.CategData{UserId: parsedUserId, Name: body.Name}
	if err = app.UpdateCategory(ctx, categId, categParams); err != nil {
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

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateFileReq](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Caso nada seja requisitado para alterar
	if body.CategId == "" && body.Name == "" && body.Extension == "" &&
		body.Mimetype == "" && len(body.Content) == 0 {
		return c.JSON(http.StatusBadRequest, BadRequestMessage)
	}

	// Transformar CategId para atualização e verificar existência
	parsedCategId := uuid.Nil
	if body.CategId != "" {
		parsedCategId, err = uuid.Parse(body.CategId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
		}
		if _, err = app.QueryCategoryById(ctx, parsedCategId); err != nil {
			return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
		}
	}

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	fileId, err := ParseEntityUUID(c, File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil && file.CategId != categId.String() {
		return c.JSON(http.StatusNotFound, FileNotFoundMessage)
	}

	// Alteração
	fileParams := app.FileData{
		CategId:   parsedCategId,
		Name:      body.Name,
		Extension: body.Extension,
		Mimetype:  body.Mimetype,
		Content:   &body.Content,
	}
	if err = app.UpdateFile(ctx, fileId, fileParams); err != nil {
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	// Remoção do usuário
	if err = app.DeleteUser(ctx, userId); err != nil {
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	// Remoção da categoria
	if err = app.DeleteCategory(ctx, categId); err != nil {
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
	// Cabeçalho e contexto da aplicação
	ctx := context.GetContext(c)
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Checar se é admin
	if admin := auth.AuthenticateAdmin(c); !admin {
		return c.JSON(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Parâmetros da URL e verificar se usuário, categoria e arquivo existem
	userId, err := ParseEntityUUID(c, User)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidUserIdMessage)
	}
	if _, err = app.QueryUserById(ctx, userId); err != nil {
		return c.JSON(http.StatusNotFound, UserNotFoundMessage)
	}

	categId, err := ParseEntityUUID(c, Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil && categ.UserId != userId.String() {
		return c.JSON(http.StatusNotFound, CategoryNotFoundMessage)
	}

	fileId, err := ParseEntityUUID(c, File)
	if err != nil {
		return c.JSON(http.StatusBadRequest, InvalidFileIdMessage)
	}
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil && file.CategId != categId.String() {
		return c.JSON(http.StatusNotFound, FileNotFoundMessage)
	}

	// Remoção do arquivo
	if err = app.DeleteFile(ctx, fileId); err != nil {
		return c.JSON(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedFileMessage)
}

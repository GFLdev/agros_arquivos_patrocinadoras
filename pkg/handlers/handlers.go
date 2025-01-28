package handlers

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/app/store"
	"agros_arquivos_patrocinadoras/pkg/auth"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// LoginHandler gerencia o processo de login.
func LoginHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[LoginReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Verificar credenciais
	userLogin, err := store.UserLogin(ctx, body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, UnauthorizedMessage)
	}
	err = bcrypt.CompareHashAndPassword([]byte(userLogin.Hash), []byte(body.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, UnauthorizedMessage)
	}

	// Gerar token e resposta
	token, err := auth.GenerateToken(c, userLogin.UserId, body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	c.Response().Header().Add(echo.HeaderAuthorization, "Bearer "+token)
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

// DownloadHandler gerencia o envio de um arquivo para download.
func DownloadHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obtenção do contexto da aplicação e dos metadados do arquivo
	ctx := context.GetContext(c)
	file, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, FileNotFoundMessage)
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
		return echo.NewHTTPError(http.StatusNotFound, FileNotFoundMessage)
	}
	c.Response().Header().Add(echo.HeaderContentType, file.Mimetype)
	return c.Attachment(path, file.Name)
}

// CreateUserHandler gerencia a criação de um usuário.
func CreateUserHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateUserReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Criptografar senha
	if len(body.Password) < 4 {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}
	hash, err := HashPassword(ctx, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	// Criar usuário
	user := store.UserParams{
		Name:     body.Name,
		Password: hash,
	}
	err = store.CreateUser(ctx, &user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, CreatedUserMessage)
}

// CreateCategoryHandler gerencia a criação de uma categoria.
func CreateCategoryHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateCategoryReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return err
	}

	// Criar categoria
	categ := store.CategParams{
		UserId: userId,
		Name:   body.Name,
	}
	err = store.CreateCategory(ctx, &categ)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, CreatedCategoryMessage)
}

// CreateFileHandler gerencia a criação de um arquivo.
func CreateFileHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[CreateFileReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Criar arquivo
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
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, CreatedFileMessage)
}

// GetAllUsers obtém todos os usuários do repositório.
func GetAllUsers(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e de todos os usuários
	ctx := context.GetContext(c)
	res, err := store.QueryAllUsers(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, UsersNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetUserById obtém um usuário com base em seu Id.
func GetUserById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obtenção do contexto da aplicação e do usuário
	ctx := context.GetContext(c)
	res, err := store.QueryUserById(ctx, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, UserNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetAllCategories obtém todas as categorias de um usuário do repositório.
func GetAllCategories(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obtenção do contexto da aplicação e de todas as categorias
	ctx := context.GetContext(c)
	res, err := store.QueryAllCategories(ctx, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, CategoriesNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetCategoryById obtém uma categoria com base em seu Id.
func GetCategoryById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	_, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obtenção do contexto da aplicação e da categoria
	ctx := context.GetContext(c)
	res, err := store.QueryCategoryById(ctx, categId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, CategoryNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetAllFiles obtém todos os arquivos de uma categoria de um usuário do
// repositório.
func GetAllFiles(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	_, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obtenção do contexto da aplicação e de todos os arquivos
	ctx := context.GetContext(c)
	res, err := store.QueryAllFiles(ctx, categId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, FilesNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// GetFileById obtém um arquivo com base em seu Id.
func GetFileById(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros
	_, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	_, err = ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obtenção do contexto da aplicação e do arquivo
	ctx := context.GetContext(c)
	res, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, UserNotFoundMessage)
	}
	return c.JSON(http.StatusOK, res)
}

// UpdateUserHandler gerencia a modificação de um usuário pelo seu Id.
func UpdateUserHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateUserReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obter os dados do usuário
	data, err := store.QueryUserById(ctx, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, UserNotFoundMessage)
	}

	// Criptografar senha, se for passado valor
	if len(body.Password) < 4 {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}
	hash, err := HashPassword(ctx, body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
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
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, UpdatedUserMessage)
}

// UpdateCategoryHandler gerencia a modificação de uma categoria pelo seu Id.
func UpdateCategoryHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateCategoryReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obter dados do usuário
	data, err := store.QueryCategoryById(ctx, categId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, CategoryNotFoundMessage)
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
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, UpdatedCategoryMessage)
}

// UpdateFileHandler gerencia a modificação de um arquivo pelo seu Id.
func UpdateFileHandler(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Obtenção do contexto da aplicação e do corpo da requisição
	ctx := context.GetContext(c)
	body, err := BodyUnmarshall[UpdateFileReq](c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, BadRequestMessage)
	}

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obter dados do usuário
	data, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, UserNotFoundMessage)
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
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}

	return c.JSON(http.StatusOK, UpdatedFileMessage)
}

// DeleteUser gerencia a exclusão de um usuário pelo seu Id.
func DeleteUser(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}

	// Obtenção do contexto da aplicação e remoção do usuário
	ctx := context.GetContext(c)
	user := store.UserDelete{UserId: userId}
	err = store.DeleteUser(ctx, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedUserMessage)
}

// DeleteCategory gerencia a exclusão de uma categoria pelo seu Id.
func DeleteCategory(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}

	// Obtenção do contexto da aplicação e remoção da categoria
	ctx := context.GetContext(c)
	categ := store.CategDelete{
		UserId:  userId,
		CategId: categId,
	}
	err = store.DeleteCategory(ctx, categ)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedCategoryMessage)
}

// DeleteFile gerencia a exclusão de um arquivo pelo seu Id.
func DeleteFile(c echo.Context) error {
	// Cabeçalho
	c.Response().Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Parâmetros da URL
	userId, err := ParseEntityUUID(c, fs.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidUserIdMessage)
	}
	categId, err := ParseEntityUUID(c, fs.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidCategoryIdMessage)
	}
	fileId, err := ParseEntityUUID(c, fs.File)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, InvalidFileIdMessage)
	}

	// Obtenção do contexto da aplicação e dos dados do arquivo
	ctx := context.GetContext(c)
	data, err := store.QueryFileById(ctx, fileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, FileNotFoundMessage)
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
		return echo.NewHTTPError(http.StatusInternalServerError, InternalServerErrorMessage)
	}
	return c.JSON(http.StatusOK, DeletedFileMessage)
}

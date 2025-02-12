package test

import (
	"agros_arquivos_patrocinadoras/pkg/app"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	h "agros_arquivos_patrocinadoras/pkg/handlers"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func echoNewContext(req *http.Request, rec *httptest.ResponseRecorder) echo.Context {
	e := echo.New()
	c := e.NewContext(req, rec)
	c.Set("appContext", newContext())
	return c
}

func TestHandlers_CreateUser(t *testing.T) {
	// Mock
	validUserJSON := `{"name": "User1", "password": "123456789"}`
	invalidUserJSON := `{"eman": "InvalidField"}`
	invalidPasswordLen := `{"name": "InvalidPassword", "password": "123"}`
	missingRequiredFields := `{}`

	// Cenário positivo
	t.Run(
		"Deve_Retornar_Created_Quando_Usuario_Criado_Com_Sucesso",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(validUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Contains(
					t,
					rec.Body.String(),
					fmt.Sprintf(`"message":"%s"`, h.CreatedUserMessage),
				)
			}

			// Verificar
			req = httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(validUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec = httptest.NewRecorder()
			c = echoNewContext(req, rec)

			if assert.NoError(t, h.LoginHandler(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(
					t,
					rec.Body.String(),
					fmt.Sprintf(`"message":"%s"`, h.LoginSuccessMessage),
				)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Bad_Request_Quando_JSON_Invalido_Recebido",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(invalidUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Conflict_Quando_Usuario_Ja_Existe",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(validUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(c)) {
				assert.Equal(t, http.StatusConflict, rec.Code)
				assert.Contains(t, rec.Body.String(), h.DuplicateUserMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Senha_Invalida",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(invalidPasswordLen),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidPasswordMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Campos_Obrigatorios_Nao_Foram_Enviados",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(missingRequiredFields),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)
}

func TestHandlers_CreateCategory(t *testing.T) {
	// Contexto
	ctx := newContext()

	// Criar usuário
	userParams := app.LoginParams{
		Username: "CategUser1",
		Password: "test123456789",
	}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	// Mock
	validCategJSON := `{"name": "Categ1"}`
	invalidCategJSON := `{"eman": "InvalidField"}`
	missingRequiredFields := `{}`

	// Cenário positivo
	t.Run(
		"Deve_Retornar_Created_Quando_Categoria_Criada_Com_Sucesso",
		func(t *testing.T) {
			// Criar categoria
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category",
				strings.NewReader(validCategJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.CreateCategoryHandler(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Contains(
					t,
					rec.Body.String(),
					fmt.Sprintf(`"message":"%s"`, h.CreatedCategoryMessage),
				)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Bad_Request_Quando_JSON_Invalido_Recebido",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category",
				strings.NewReader(invalidCategJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.CreateCategoryHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Existe",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+invalidUserId+"/category",
				strings.NewReader(validCategJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category")
			c.SetParamNames("userId")
			c.SetParamValues(invalidUserId)

			if assert.NoError(t, h.CreateCategoryHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Campos_Obrigatorios_Nao_Foram_Enviados",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category",
				strings.NewReader(missingRequiredFields),
			)
			req.SetPathValue("userId", userId.String())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.CreateCategoryHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)
}

func TestHandlers_CreateFile(t *testing.T) {
	// Contexto
	ctx := newContext()

	// Criar usuário e categoria
	userParams := app.LoginParams{
		Username: "FileUser1",
		Password: "123456789",
	}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)
	categParams := app.CategData{
		UserId: userId,
		Name:   "FileCateg1",
	}
	categId, err := app.CreateCategory(ctx, categParams)
	assert.NoError(t, err)

	// Mock
	validFileJSON := `{
		"name": "File1",
		"extension": ".txt",
		"mimetype": "text/plain",
		"content": "SGVsbG8gV29ybGQ="
	}`
	invalidFileJSON := `{"eman": "InvalidField"}`
	missingRequiredFields := `{"name": "MissingFields"}`

	// Cenário positivo
	t.Run(
		"Deve_Retornar_Created_Quando_Arquivo_Criado_Com_Sucesso",
		func(t *testing.T) {
			// Criar categoria
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category"+categId.String()+"/file",
				strings.NewReader(validFileJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), categId.String())

			if assert.NoError(t, h.CreateFileHandler(c)) {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Contains(
					t,
					rec.Body.String(),
					fmt.Sprintf(`"message":"%s"`, h.CreatedFileMessage),
				)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Bad_Request_Quando_JSON_Invalido_Recebido",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category"+categId.String()+"/file",
				strings.NewReader(invalidFileJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), categId.String())

			if assert.NoError(t, h.CreateFileHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Existe",
		func(t *testing.T) {
			invalidCategId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category"+invalidCategId+"/file",
				strings.NewReader(validFileJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), invalidCategId)

			if assert.NoError(t, h.CreateFileHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Campos_Obrigatorios_Nao_Foram_Enviados",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user/"+userId.String()+"/category"+categId.String()+"/file",
				strings.NewReader(missingRequiredFields),
			)
			req.SetPathValue("userId", userId.String())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), categId.String())

			if assert.NoError(t, h.CreateFileHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)
}

func TestHandlers_ReadUser(t *testing.T) {
	// Mock
	ctx := newContext()
	user1Params := app.LoginParams{Username: "ReadUser1", Password: "123456789"}
	user2Params := app.LoginParams{Username: "ReadUser2", Password: "123456789"}
	user1Id, err := app.CreateUser(ctx, user1Params)
	assert.NoError(t, err)
	user2Id, err := app.CreateUser(ctx, user2Params)
	assert.NoError(t, err)

	// Cenários positivos
	t.Run(
		"Deve_Retornar_OK_Quando_Usuario_Encontrado",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+user1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(user1Id.String())

			if assert.NoError(t, h.GetUserById(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), `"user_id":"`+user1Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+user1Params.Username+`"`)
				assert.Contains(t, rec.Body.String(), `"password":""`)
			}
		},
	)

	t.Run(
		"Deve_Retornar_OK_Quando_Usuarios_Encontrados",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/user",
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user")

			if assert.NoError(t, h.GetAllUsers(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), `"user_id":"`+user1Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+user1Params.Username+`"`)
				assert.Contains(t, rec.Body.String(), `"user_id":"`+user2Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+user2Params.Username+`"`)
				assert.Contains(t, rec.Body.String(), `"password":""`)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Encontrado",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+invalidUserId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(invalidUserId)

			if assert.NoError(t, h.GetUserById(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Usuario_Invalido",
		func(t *testing.T) {
			invalidUserId := "InvalidUserId"
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+invalidUserId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(invalidUserId)

			if assert.NoError(t, h.GetUserById(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidUserIdMessage)
			}
		},
	)
}

func TestHandlers_ReadCategory(t *testing.T) {
	// Mock
	ctx := newContext()
	userParams := app.LoginParams{Username: "ReadCategUser", Password: "123456789"}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	categ1Params := app.CategData{UserId: userId, Name: "ReadCateg1"}
	categ2Params := app.CategData{UserId: userId, Name: "ReadCateg2"}
	categ1Id, err := app.CreateCategory(ctx, categ1Params)
	assert.NoError(t, err)
	categ2Id, err := app.CreateCategory(ctx, categ2Params)
	assert.NoError(t, err)

	// Cenários positivos
	t.Run(
		"Deve_Retornar_OK_Quando_Categoria_Encontrada",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+categ1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), categ1Id.String())

			if assert.NoError(t, h.GetCategoryById(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), `"categ_id":"`+categ1Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"user_id":"`+userId.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+categ1Params.Name+`"`)
			}
		},
	)

	t.Run(
		"Deve_Retornar_OK_Quando_Categorias_Encontradas",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category",
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.GetAllCategories(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), `"categ_id":"`+categ1Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+categ1Params.Name+`"`)
				assert.Contains(t, rec.Body.String(), `"categ_id":"`+categ2Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+categ2Params.Name+`"`)
				assert.Contains(t, rec.Body.String(), `"user_id":"`+userId.String()+`"`)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Encontrada",
		func(t *testing.T) {
			invalidCategId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+invalidCategId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), invalidCategId)

			if assert.NoError(t, h.GetCategoryById(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Encontrado",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+invalidUserId+"/category/"+categ1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(invalidUserId, categ1Id.String())

			if assert.NoError(t, h.GetCategoryById(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Categoria_Invalida",
		func(t *testing.T) {
			invalidCategId := "InvalidCategId"
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+invalidCategId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), invalidCategId)

			if assert.NoError(t, h.GetCategoryById(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidCategoryIdMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Usuario_Invalido",
		func(t *testing.T) {
			invalidUserId := "InvalidUserId"
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+invalidUserId+"/category/"+categ1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(invalidUserId, categ1Id.String())

			if assert.NoError(t, h.GetCategoryById(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidUserIdMessage)
			}
		},
	)
}

func TestHandlers_ReadFile(t *testing.T) {
	// Mock
	ctx := newContext()
	userParams := app.LoginParams{Username: "ReadFileUser", Password: "123456789"}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	categParams := app.CategData{UserId: userId, Name: "ReadFileCateg"}
	categId, err := app.CreateCategory(ctx, categParams)
	assert.NoError(t, err)

	content := []byte("Test")
	file1Params := app.FileData{
		CategId:   categId,
		Name:      "ReadFile1",
		Extension: ".txt",
		Mimetype:  "text/plain",
		Content:   &content,
	}
	file2Params := app.FileData{
		CategId:   categId,
		Name:      "ReadFile2",
		Extension: ".txt",
		Mimetype:  "text/plain",
		Content:   &content,
	}
	file1Id, err := app.CreateFile(ctx, file1Params)
	assert.NoError(t, err)
	file2Id, err := app.CreateFile(ctx, file2Params)
	assert.NoError(t, err)

	// Cenários positivos
	t.Run(
		"Deve_Retornar_OK_Quando_Arquivo_Encontrado",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file/"+file1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categId.String(), file1Id.String())

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Header().Get(echo.HeaderContentType), file1Params.Mimetype)
				assert.Contains(t, rec.Body.String(), string(*file1Params.Content))
			}
		},
	)

	t.Run(
		"Deve_Retornar_OK_Quando_Arquivos_Encontrados",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file",
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), categId.String())

			if assert.NoError(t, h.GetAllFiles(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), `"file_id":"`+file1Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+file1Params.Name+`"`)
				assert.Contains(t, rec.Body.String(), `"extension":"`+file1Params.Extension+`"`)
				assert.Contains(t, rec.Body.String(), `"mimetype":"`+file1Params.Mimetype+`"`)
				assert.Contains(t, rec.Body.String(), `"file_id":"`+file2Id.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"name":"`+file2Params.Name+`"`)
				assert.Contains(t, rec.Body.String(), `"extension":"`+file2Params.Extension+`"`)
				assert.Contains(t, rec.Body.String(), `"mimetype":"`+file2Params.Mimetype+`"`)
				assert.Contains(t, rec.Body.String(), `"categ_id":"`+categId.String()+`"`)
				assert.Contains(t, rec.Body.String(), `"blob":null`)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Not_Found_Quando_Arquivo_Nao_Encontrado",
		func(t *testing.T) {
			invalidFileId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file/"+invalidFileId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categId.String(), invalidFileId)

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.FileNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Encontrada",
		func(t *testing.T) {
			invalidCategId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+invalidCategId+"/file/"+file1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), invalidCategId, file1Id.String())

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Encontrado",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+invalidUserId+"/category/"+categId.String()+"/file/"+file1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(invalidUserId, categId.String(), file1Id.String())

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Arquivo_Invalido",
		func(t *testing.T) {
			invalidFileId := "InvalidFileId"
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file/"+invalidFileId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categId.String(), invalidFileId)

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidFileIdMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Categoria_Invalida",
		func(t *testing.T) {
			invalidCategId := "InvalidCategId"
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+userId.String()+"/category/"+invalidCategId+"/file/"+file1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), invalidCategId, file1Id.String())

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidCategoryIdMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Usuario_Invalido",
		func(t *testing.T) {
			invalidUserId := "InvalidUserId"
			req := httptest.NewRequest(
				http.MethodGet,
				"/user/"+invalidUserId+"/category/"+categId.String()+"/file/"+file1Id.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(invalidUserId, categId.String(), file1Id.String())

			if assert.NoError(t, h.GetFileById(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidUserIdMessage)
			}
		},
	)
}

func updateCategTestHelper(
	t *testing.T,
	ctx *context.Context,
	ids [2]uuid.UUID,
	payload string,
	categParams app.CategData,
) {
	req := httptest.NewRequest(
		http.MethodPatch,
		"/user/"+ids[0].String()+"/category/"+ids[1].String(),
		strings.NewReader(payload),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echoNewContext(req, rec)
	c.SetPath("/user/:userId/category/:categId")
	c.SetParamNames("userId", "categId")
	c.SetParamValues(ids[0].String(), ids[1].String())

	if assert.NoError(t, h.UpdateCategoryHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), h.UpdatedCategoryMessage)
	}
	categ, err := app.QueryCategoryById(ctx, ids[1])
	if assert.NoError(t, err) {
		assert.Equal(t, categParams.UserId.String(), categ.UserId)
		assert.Equal(t, categParams.Name, categ.Name)
	}
}

func updateFileTestHelper(
	t *testing.T,
	ctx *context.Context,
	ids [3]uuid.UUID,
	payload string,
	fileParams app.FileData,
) {
	req := httptest.NewRequest(
		http.MethodPatch,
		"/user/"+ids[0].String()+"/category/"+ids[1].String()+"/file/"+ids[2].String(),
		strings.NewReader(payload),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echoNewContext(req, rec)
	c.SetPath("/user/:userId/category/:categId/file/:fileId")
	c.SetParamNames("userId", "categId", "fileId")
	c.SetParamValues(ids[0].String(), ids[1].String(), ids[2].String())

	if assert.NoError(t, h.UpdateFileHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), h.UpdatedFileMessage)
	}
	file, err := app.QueryFileById(ctx, ids[2])
	if assert.NoError(t, err) {
		assert.Equal(t, fileParams.CategId.String(), file.CategId)
		assert.Equal(t, fileParams.Name, file.Name)
		assert.Equal(t, fileParams.Extension, file.Extension)
		assert.Equal(t, fileParams.Mimetype, file.Mimetype)
		assert.Equal(t, *fileParams.Content, file.Blob)
	}
}

func TestHandlers_UpdateUser(t *testing.T) {
	// Contexto
	ctx := newContext()

	// Criar usuário
	userParams := app.LoginParams{
		Username: "BeforeUser1",
		Password: "123456789",
	}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	// Cenário positivo
	t.Run(
		"Deve_Retornar_OK_Quando_Usuario_Atualizado_Com_Sucesso",
		func(t *testing.T) {
			// Mock
			newName := "AfterUser1"
			newPwd := "987654321"

			// Atualizar nome
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String(),
				strings.NewReader(`{"name":"`+newName+`"}`),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.UpdateUserHandler(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UpdatedUserMessage)
			}

			// Verificar os dados
			newParams := app.LoginParams{Username: newName, Password: userParams.Password}
			id, err := app.QueryLogin(ctx, newParams)
			if assert.NoError(t, err) {
				assert.Equal(t, userId, id)
			}

			// Atualizar senha
			req = httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String(),
				strings.NewReader(`{"password":"`+newPwd+`"}`),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec = httptest.NewRecorder()
			c = echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.UpdateUserHandler(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UpdatedUserMessage)
			}

			// Verificar os dados
			newParams = app.LoginParams{Username: newName, Password: newPwd}
			id, err = app.QueryLogin(ctx, newParams)
			if assert.NoError(t, err) {
				assert.Equal(t, userId, id)
			}
		},
	)

	// Cenário negativo
	t.Run(
		"Deve_Retornar_Bad_Request_Quando_JSON_Invalido_Recebido",
		func(t *testing.T) {
			invalidJSON := `{"eman":"InvalidField"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String(),
				strings.NewReader(invalidJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.UpdateUserHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Campos_Obrigatorios_Nao_Foram_Enviados",
		func(t *testing.T) {
			missingRequiredFields := `{}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String(),
				strings.NewReader(missingRequiredFields),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.UpdateUserHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Existe",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			validUpdateJSON := `{"name": "UpdatedName"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+invalidUserId,
				strings.NewReader(validUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(invalidUserId)

			if assert.NoError(t, h.UpdateUserHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Senha_For_Muito_Curta",
		func(t *testing.T) {
			invalidPasswordJSON := `{"password": "123"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String(),
				strings.NewReader(invalidPasswordJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.UpdateUserHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidPasswordMessage)
			}
		},
	)
}

func TestHandlers_UpdateCategory(t *testing.T) {
	// Contexto
	ctx := newContext()

	// Criar usuários
	user1Params := app.LoginParams{Username: "UpCategUser1", Password: "123456789"}
	user2Params := app.LoginParams{Username: "UpCategUser2", Password: "123456789"}

	user1Id, err := app.CreateUser(ctx, user1Params)
	assert.NoError(t, err)
	user2Id, err := app.CreateUser(ctx, user2Params)
	assert.NoError(t, err)

	// Criar categoria
	categ1Params := app.CategData{UserId: user1Id, Name: "BeforeCategory1"}
	categ2Params := app.CategData{UserId: user2Id, Name: "BeforeCategory2"}

	categ1Id, err := app.CreateCategory(ctx, categ1Params)
	assert.NoError(t, err)
	categ2Id, err := app.CreateCategory(ctx, categ2Params)
	assert.NoError(t, err)

	// Cenário positivo
	t.Run(
		"Deve_Retornar_OK_Quando_Categoria_Atualizada_Com_Sucesso",
		func(t *testing.T) {
			// Mock
			newName := "AfterCategory1"
			ids := [2]uuid.UUID{user1Id, categ1Id}

			// Atualizar nome e verificar dados
			categ1Params.Name = newName
			payload := `{"name": "` + newName + `"}`
			updateCategTestHelper(t, ctx, ids, payload, categ1Params)

			// Atualizar id de usuário
			categ1Params.UserId = user2Id
			payload = `{"userId": "` + user2Id.String() + `"}`
			updateCategTestHelper(t, ctx, ids, payload, categ1Params)
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Bad_Request_Quando_JSON_Invalido_Recebido",
		func(t *testing.T) {
			invalidUpdateJSON := `{"eman": "InvalidField"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+user2Id.String()+"/category/"+categ2Id.String(),
				strings.NewReader(invalidUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(user2Id.String(), categ2Id.String())

			if assert.NoError(t, h.UpdateCategoryHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Campos_Obrigatorios_Ausentes",
		func(t *testing.T) {
			missingFieldsJSON := `{}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+user2Id.String()+"/category/"+categ2Id.String(),
				strings.NewReader(missingFieldsJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(user2Id.String(), categ2Id.String())

			if assert.NoError(t, h.UpdateCategoryHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Existe",
		func(t *testing.T) {
			invalidCategoryId := uuid.New().String()
			validUpdateJSON := `{"name": "UpdatedName"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+user2Id.String()+"/category/"+invalidCategoryId,
				strings.NewReader(validUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(user2Id.String(), invalidCategoryId)

			if assert.NoError(t, h.UpdateCategoryHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Existe",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			validUpdateJSON := `{"name": "UpdatedName"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+invalidUserId+"/category/"+categ2Id.String(),
				strings.NewReader(validUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(invalidUserId, categ2Id.String())

			if assert.NoError(t, h.UpdateCategoryHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)
}

func TestHandlers_UpdateFile(t *testing.T) {
	// Contexto
	ctx := newContext()

	// Criar usuário e mock JSON
	validUpdateJSON := `{"name": "UpdatedName"}`
	userParams := app.LoginParams{Username: "UpFileUser1", Password: "123456789"}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	// Criar categorias
	categ1Params := app.CategData{UserId: userId, Name: "UpFileCateg1"}
	categ2Params := app.CategData{UserId: userId, Name: "UpFileCateg2"}

	categ1Id, err := app.CreateCategory(ctx, categ1Params)
	assert.NoError(t, err)
	categ2Id, err := app.CreateCategory(ctx, categ2Params)
	assert.NoError(t, err)

	// Criar arquivos
	originalContent := []byte("Original Content")
	file1Params := app.FileData{
		CategId:   categ1Id,
		Name:      "BeforeFile1",
		Extension: ".txt",
		Mimetype:  "text/plain",
		Content:   &originalContent,
	}
	file2Params := app.FileData{
		CategId:   categ2Id,
		Name:      "BeforeFile2",
		Extension: ".txt",
		Mimetype:  "text/plain",
		Content:   &originalContent,
	}

	file1Id, err := app.CreateFile(ctx, file1Params)
	assert.NoError(t, err)
	file2Id, err := app.CreateFile(ctx, file2Params)
	assert.NoError(t, err)

	// Cenário positivo
	t.Run(
		"Deve_Retornar_OK_Quando_Arquivo_Atualizado_Com_Sucesso",
		func(t *testing.T) {
			// Mock
			newName := "AfterFile1"
			newExtension := ".csv"
			newMimetype := "text/csv"
			newContent := []byte("a,b,c")
			ids := [3]uuid.UUID{userId, categ1Id, file1Id}

			// Atualizar nome e verificar dados
			file1Params.Name = newName
			payload := `{"name": "` + newName + `"}`
			updateFileTestHelper(t, ctx, ids, payload, file1Params)

			// Atualizar extensão e verificar dados
			file1Params.Extension = newExtension
			payload = `{"extension": "` + newExtension + `"}`
			updateFileTestHelper(t, ctx, ids, payload, file1Params)

			// Atualizar mimetype e verificar dados
			file1Params.Mimetype = newMimetype
			payload = `{"mimetype": "` + newMimetype + `"}`
			updateFileTestHelper(t, ctx, ids, payload, file1Params)

			// Atualizar conteúdo e verificar dados
			encodedContent := base64.StdEncoding.EncodeToString(newContent)
			file1Params.Content = &newContent
			payload = `{"content": "` + encodedContent + `"}`
			updateFileTestHelper(t, ctx, ids, payload, file1Params)

			// Atualizar CategId e verificar dados
			file1Params.CategId = categ2Id
			payload = `{"categId": "` + categ2Id.String() + `"}`
			updateFileTestHelper(t, ctx, ids, payload, file1Params)
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Bad_Request_Quando_JSON_Invalido_Recebido",
		func(t *testing.T) {
			invalidUpdateJSON := `{"eman": "InvalidField"}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String()+"/category/"+categ2Id.String()+"/file/"+file2Id.String(),
				strings.NewReader(invalidUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categ2Id.String(), file2Id.String())

			if assert.NoError(t, h.UpdateFileHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Campos_Obrigatorios_Ausentes",
		func(t *testing.T) {
			missingFieldsJSON := `{}`
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String()+"/category/"+categ2Id.String()+"/file/"+file2Id.String(),
				strings.NewReader(missingFieldsJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categ2Id.String(), file2Id.String())

			if assert.NoError(t, h.UpdateFileHandler(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Arquivo_Nao_Existe",
		func(t *testing.T) {
			invalidFileId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String()+"/category/"+categ1Id.String()+"/file/"+invalidFileId,
				strings.NewReader(validUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categ2Id.String(), invalidFileId)

			if assert.NoError(t, h.UpdateFileHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.FileNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Existe",
		func(t *testing.T) {
			invalidCategId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+userId.String()+"/category/"+invalidCategId+"/file/"+file2Id.String(),
				strings.NewReader(validUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), invalidCategId, file2Id.String())

			if assert.NoError(t, h.UpdateFileHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Existe",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodPatch,
				"/user/"+invalidUserId+"/category/"+categ2Id.String()+"/file/"+file2Id.String(),
				strings.NewReader(validUpdateJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(invalidUserId, categ2Id.String(), file2Id.String())

			if assert.NoError(t, h.UpdateFileHandler(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)
}

func TestHandlers_DeleteUser(t *testing.T) {
	// Mock
	ctx := newContext()
	userParams := app.LoginParams{Username: "DeleteUser", Password: "123456789"}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	// Cenários positivos
	t.Run(
		"Deve_Retornar_OK_Quando_Usuario_Excluido_Com_Sucesso",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(userId.String())

			if assert.NoError(t, h.DeleteUser(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), h.DeletedUserMessage)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Encontrado",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+invalidUserId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(invalidUserId)

			if assert.NoError(t, h.DeleteUser(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Usuario_Invalido",
		func(t *testing.T) {
			invalidUserId := "InvalidUserId"
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+invalidUserId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId")
			c.SetParamNames("userId")
			c.SetParamValues(invalidUserId)

			if assert.NoError(t, h.DeleteUser(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidUserIdMessage)
			}
		},
	)
}

func TestHandlers_DeleteCategory(t *testing.T) {
	// Mock
	ctx := newContext()
	userParams := app.LoginParams{Username: "DeleteCategUser", Password: "123456789"}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	categParams := app.CategData{UserId: userId, Name: "DeleteCateg"}
	categId, err := app.CreateCategory(ctx, categParams)
	assert.NoError(t, err)

	// Cenários positivos
	t.Run(
		"Deve_Retornar_OK_Quando_Categoria_Deletada_Com_Sucesso",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+categId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), categId.String())

			if assert.NoError(t, h.DeleteCategory(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), h.DeletedCategoryMessage)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Encontrada",
		func(t *testing.T) {
			invalidCategId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+invalidCategId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), invalidCategId)

			if assert.NoError(t, h.DeleteCategory(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Encontrado",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+invalidUserId+"/category/"+categId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(invalidUserId, categId.String())

			if assert.NoError(t, h.DeleteCategory(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Categoria_Invalida",
		func(t *testing.T) {
			invalidCategId := "InvalidCategId"
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+invalidCategId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(userId.String(), invalidCategId)

			if assert.NoError(t, h.DeleteCategory(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidCategoryIdMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Usuario_Invalido",
		func(t *testing.T) {
			invalidUserId := "InvalidUserId"
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+invalidUserId+"/category/"+categId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId")
			c.SetParamNames("userId", "categId")
			c.SetParamValues(invalidUserId, categId.String())

			if assert.NoError(t, h.DeleteCategory(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidUserIdMessage)
			}
		},
	)
}

func TestHandlers_DeleteFile(t *testing.T) {
	// Mock
	ctx := newContext()
	userParams := app.LoginParams{Username: "DeleteFileUser", Password: "123456789"}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	categParams := app.CategData{UserId: userId, Name: "DeleteFileCateg"}
	categId, err := app.CreateCategory(ctx, categParams)
	assert.NoError(t, err)

	content := []byte("Test")
	fileParams := app.FileData{
		CategId:   categId,
		Name:      "ReadFile1",
		Extension: ".txt",
		Mimetype:  "text/plain",
		Content:   &content,
	}
	fileId, err := app.CreateFile(ctx, fileParams)
	assert.NoError(t, err)

	// Cenários positivos
	t.Run(
		"Deve_Retornar_OK_Quando_Arquivo_Encontrado",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file/"+fileId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categId.String(), fileId.String())

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Contains(t, rec.Body.String(), h.DeletedFileMessage)
			}
		},
	)

	// Cenários negativos
	t.Run(
		"Deve_Retornar_Not_Found_Quando_Arquivo_Nao_Encontrado",
		func(t *testing.T) {
			invalidFileId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file/"+invalidFileId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categId.String(), invalidFileId)

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.FileNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Categoria_Nao_Encontrada",
		func(t *testing.T) {
			invalidCategId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+invalidCategId+"/file/"+fileId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), invalidCategId, fileId.String())

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.CategoryNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Not_Found_Quando_Usuario_Nao_Encontrado",
		func(t *testing.T) {
			invalidUserId := uuid.New().String()
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+invalidUserId+"/category/"+categId.String()+"/file/"+fileId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(invalidUserId, categId.String(), fileId.String())

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusNotFound, rec.Code)
				assert.Contains(t, rec.Body.String(), h.UserNotFoundMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Arquivo_Invalido",
		func(t *testing.T) {
			invalidFileId := "InvalidFileId"
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+categId.String()+"/file/"+invalidFileId,
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), categId.String(), invalidFileId)

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidFileIdMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Categoria_Invalida",
		func(t *testing.T) {
			invalidCategId := "InvalidCategId"
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+userId.String()+"/category/"+invalidCategId+"/file/"+fileId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(userId.String(), invalidCategId, fileId.String())

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidCategoryIdMessage)
			}
		},
	)

	t.Run(
		"Deve_Retornar_Bad_Request_Quando_Usuario_Invalido",
		func(t *testing.T) {
			invalidUserId := "InvalidUserId"
			req := httptest.NewRequest(
				http.MethodDelete,
				"/user/"+invalidUserId+"/category/"+categId.String()+"/file/"+fileId.String(),
				nil,
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)
			c.SetPath("/user/:userId/category/:categId/file/:fileId")
			c.SetParamNames("userId", "categId", "fileId")
			c.SetParamValues(invalidUserId, categId.String(), fileId.String())

			if assert.NoError(t, h.DeleteFile(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidUserIdMessage)
			}
		},
	)
}

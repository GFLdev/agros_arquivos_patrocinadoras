package test

import (
	"agros_arquivos_patrocinadoras/pkg/app"
	h "agros_arquivos_patrocinadoras/pkg/handlers"
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
	appCtx := newContext()
	echoCtx := e.NewContext(req, rec)
	echoCtx.Set("appContext", appCtx)
	return echoCtx
}

func TestHandlers_CreateUser(t *testing.T) {
	// Mock
	validUserJSON := `{"name": "testUser1", "password": "test123456789"}`
	invalidUserJSON := `{"eman": "testUser2"}`
	invalidPasswordLen := `{"name": "testUser3", "password": "123"}`
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
	// Criar usuário
	ctx := newContext()
	userParams := app.UserParams{
		Name:     "testCateg1",
		Password: "test123456789",
	}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)

	// Mock
	validCategJSON := `{"name": "test1"}`
	invalidCategJSON := `{"eman": "test2"}`
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
		"Deve_Retornar_Internal_Server_Error_Quando_Usuario_Nao_Existe",
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
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InternalServerErrorMessage)
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
	// Criar usuário e categoria
	ctx := newContext()
	userParams := app.UserParams{
		Name:     "testFile1",
		Password: "test123456789",
	}
	userId, err := app.CreateUser(ctx, userParams)
	assert.NoError(t, err)
	categParams := app.CategParams{
		UserId: userId,
		Name:   "test",
	}
	categId, err := app.CreateCategory(ctx, categParams)
	assert.NoError(t, err)

	// Mock
	validFileJSON := `{
		"name": "test1",
		"extension": ".txt",
		"mimetype": "text/plain",
		"content": "SGVsbG8gV29ybGQ="
	}`
	invalidFileJSON := `{"eman": "test2"}`
	missingRequiredFields := `{"name": "test1"}`

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
		"Deve_Retornar_Internal_Server_Error_Quando_Categoria_Nao_Existe",
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
				assert.Equal(t, http.StatusInternalServerError, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InternalServerErrorMessage)
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

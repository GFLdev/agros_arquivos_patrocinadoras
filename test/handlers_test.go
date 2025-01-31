package test

import (
	h "agros_arquivos_patrocinadoras/pkg/handlers"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func echoNewContext(req *http.Request, rec *httptest.ResponseRecorder) *echo.Context {
	e := echo.New()
	appCtx := newContext()
	echoCtx := e.NewContext(req, rec)
	echoCtx.Set("appContext", appCtx)
	return &echoCtx
}

func TestHandlers_CreateUser(t *testing.T) {
	// Mock
	validUserJSON := `{"name": "test1", "password": "test123456789"}`
	invalidUserJSON := `{"eman": "test2"}`
	invalidPasswordLen := `{"name": "test3", "password": "123"}`
	missingRequiredFields := `{}`

	// Cenário positivo
	t.Run(
		"Should_Return_Created_When_User_Is_Created_Successfully",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(validUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(*c)) {
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
		"Should_Return_Bad_Request_When_Body_Is_Invalid",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(invalidUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(*c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)

	t.Run(
		"Should_Return_Conflict_When_User_Already_Exists",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(validUserJSON),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(*c)) {
				assert.Equal(t, http.StatusConflict, rec.Code)
				assert.Contains(t, rec.Body.String(), h.DuplicateUserMessage)
			}
		},
	)

	t.Run(
		"Should_Return_Bad_Request_When_Password_Is_Invalid_Length",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(invalidPasswordLen),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(*c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.InvalidPasswordMessage)
			}
		},
	)

	t.Run(
		"Should_Return_Bad_Request_When_Missing_Required_Fields",
		func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/user",
				strings.NewReader(missingRequiredFields),
			)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echoNewContext(req, rec)

			if assert.NoError(t, h.CreateUserHandler(*c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				assert.Contains(t, rec.Body.String(), h.BadRequestMessage)
			}
		},
	)
}

package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/auth"
	"agros_arquivos_patrocinadoras/pkg/handlers"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ConfigRoutes define as rotas da aplicação utilizando o framework Echo.
//
// Parâmetros:
//   - e: instância do Echo, usada para registrar as rotas.
//   - ctx: contexto da aplicação que contém informações e dependências.
func ConfigRoutes(e *echo.Echo, ctx *context.Context) {
	ctx.Logger.Info("Configurando rotas")

	// Grupo para autenticação
	authGroup := e.Group("/auth")
	authGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(ctx.Config.JwtSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.CustomClaims)
		},
	}))

	// Login
	e.POST("/login", handlers.LoginHandler)

	// Sessão
	authGroup.GET("/session", handlers.SessionHandler)

	// Usuário
	authGroup.POST("/user", handlers.CreateUserHandler)
	authGroup.GET("/user", handlers.GetAllUsers)
	authGroup.GET("/user/:userId", handlers.GetUserById)
	authGroup.PATCH("/user/:userId", handlers.UpdateUserHandler)
	authGroup.DELETE("/user/:userId", handlers.DeleteUser)

	// Categorias
	authGroup.POST("/user/:userId/category", handlers.CreateCategoryHandler)
	authGroup.GET("/user/:userId/category", handlers.GetAllCategories)
	authGroup.GET("/user/:userId/category/:categId", handlers.GetCategoryById)
	authGroup.PATCH("/user/:userId/category/:categId", handlers.UpdateCategoryHandler)
	authGroup.DELETE("/user/:userId/category/:categId", handlers.DeleteCategory)

	// Arquivos
	authGroup.POST("/user/:userId/category/:categId/file", handlers.CreateFileHandler)
	authGroup.GET("/user/:userId/category/:categId/file", handlers.GetAllFiles)
	authGroup.GET("/user/:userId/category/:categId/file/:fileId", handlers.GetFileById)
	authGroup.PATCH("/user/:userId/category/:categId/file/:fileId", handlers.UpdateFileHandler)
	authGroup.DELETE("/user/:userId/category/:categId/file/:fileId", handlers.DeleteFile)

	// Preflight: rota coringa
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Rota de fallback para garantir que index.html seja carregado para todas as requisições
	e.GET("/*", func(c echo.Context) error {
		return c.File("web/dist/index.html")
	})

}

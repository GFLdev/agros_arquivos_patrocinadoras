package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/handlers"
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

	// Login
	e.POST("/login", handlers.LoginHandler)

	// Usuário
	e.POST("/user", handlers.CreateUserHandler)
	e.GET("/user", handlers.GetAllUsers)
	e.GET("/user/:userId", handlers.GetUserById)
	e.PATCH("/user/:userId", handlers.UpdateUserHandler)
	e.DELETE("/user/:userId", handlers.DeleteUser)

	// Categorias
	e.POST("/user/:userId/category", handlers.CreateCategoryHandler)
	e.GET("/user/:userId/category", handlers.GetAllCategories)
	e.GET("/user/:userId/category/:categId", handlers.GetCategoryById)
	e.PATCH("/user/:userId/category/:categId", handlers.UpdateCategoryHandler)
	e.DELETE("/user/:userId/category/:categId", handlers.DeleteCategory)

	// Arquivos
	e.POST("/user/:userId/category/:categId/file", handlers.CreateFileHandler)
	e.GET("/user/:userId/category/:categId/file", handlers.GetAllFiles)
	e.GET("/user/:userId/category/:categId/file/:fileId", handlers.GetFileById)
	e.PATCH("/user/:userId/category/:categId/file/:fileId", handlers.UpdateFileHandler)
	e.DELETE("/user/:userId/category/:categId/file/:fileId", handlers.DeleteFile)

	// Preflight: rota coringa
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}

package main

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/pkg/handlers"
	"github.com/labstack/echo/v4"
)

// ConfigRoutes define as rotas principais da aplicação, separando as
// funcionalidades para usuários e administradores. As rotas de usuários
// envolvem operações relacionadas a categorias e arquivos, enquanto as
// rotas de administradores oferecem funcionalidades adicionais para
// manipulação de usuários, categorias de arquivos e arquivos.
//
// As rotas de usuários incluem:
// - GET e POST para categorias e arquivos.
//
// As rotas de administradores incluem:
// - GET, POST, PATCH e DELETE para usuários, categorias de arquivos e arquivos.
func ConfigRoutes(e *echo.Echo, ctx *services.AppWrapper) {
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

	// Download
	e.GET("/user/:userId/category/:categId/file/:fileId/download",
		handlers.DownloadHandler)
}

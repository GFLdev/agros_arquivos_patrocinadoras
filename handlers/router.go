package handlers

import (
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
func ConfigRoutes(e *echo.Echo, ctx *AppContext) {
	ctx.Logger.Info("Configurando rotas")

	// Login
	e.POST("/login", LoginHandler)

	// Usuário
	e.POST("/user", CreateUserHandler)
	e.GET("/user", AllUsersHandler)
	e.GET("/user/:userId", UserByIdHandler)
	e.PATCH("/user/:userId", UpdateUserHandler)
	e.DELETE("/user/:userId", DeleteUserHandler)

	// Categorias
	e.POST("/user/:userId/category", CreateCategoryHandler)
	e.GET("/user/:userId/category", AllCategoriesHandler)
	e.GET("/user/:userId/category/:categId", CategoryByIdHandler)
	e.PATCH("/user/:userId/category/:categId", UpdateCategoryHandler)
	e.DELETE("/user/:userId/category/:categId", DeleteCategoryHandler)

	// Arquivos
	e.POST("/user/:userId/category/:categId/file", CreateFileHandler)
	e.GET("/user/:userId/category/:categId/file", AllFilesHandler)
	e.GET("/user/:userId/category/:categId/file/:fileId", FileByIdHandler)
	e.PATCH("/user/:userId/category/:categId/file/:fileId", UpdateFileHandler)
	e.DELETE("/user/:userId/category/:categId/file/:fileId", DeleteFileHandler)

	// Download
	e.GET("/user/:userId/category/:categId/file/:fileId/download",
		DownloadHandler)
}

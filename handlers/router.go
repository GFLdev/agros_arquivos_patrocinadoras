package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
	// Usuário - Categorias
	e.GET("/category", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.POST("/category", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	// Usuário - Arquivos
	e.GET("/category/file", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.POST("/category/file", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})

	// Administrador
	// Administrador - Usuários
	e.GET("/admin/user", AllUsersHandler)
	e.POST("/admin/user/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.PATCH("/admin/user/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.DELETE("/admin/user/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	// Administrador - Categorias de arquivos
	e.GET("/admin/category", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.GET("/admin/category/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.POST("/admin/category/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.PATCH("/admin/category/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.DELETE("/admin/category/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	// Administrador - Arquivos
	e.GET("/admin/file", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.GET("/admin/file/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.POST("/admin/file/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.PATCH("/admin/file/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	e.DELETE("/admin/file/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
}

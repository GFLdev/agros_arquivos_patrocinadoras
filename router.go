package main

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
	userGroup := e.Group("/user/:id")
	// Usuário - Categorias
	userCategoryGroup := userGroup.Group("/category")
	userCategoryGroup.GET("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	userCategoryGroup.POST("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	// Usuário - Arquivos
	userFileGroup := userGroup.Group("/file")
	userFileGroup.GET("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	userFileGroup.POST("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})

	// Administrador
	adminGroup := e.Group("/admin")
	// Administrador - Usuários
	adminUserGroup := adminGroup.Group("/user")
	adminUserGroup.GET("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminUserGroup.POST("/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminUserGroup.PATCH("/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminUserGroup.DELETE("/:id", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	// Administrador - Categorias de arquivos
	adminCategoryGroup := adminGroup.Group("/category")
	adminCategoryGroup.GET("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminCategoryGroup.POST("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminCategoryGroup.PATCH("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminCategoryGroup.DELETE("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	// Administrador - Arquivos
	adminFileGroup := adminCategoryGroup.Group("/file")
	adminFileGroup.GET("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminFileGroup.POST("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminFileGroup.PATCH("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
	adminFileGroup.DELETE("/", func(c echo.Context) error {
		c.Response().Header().Add("Content-Type", "application/json")
		return c.JSON(http.StatusOK, nil)
	})
}

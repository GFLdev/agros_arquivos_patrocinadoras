package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/handlers"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"net/http"
)

// ContextMiddleware é o middleware para implementar context.Context como
// contexto padrão a ser usado por echo.Echo.
func ContextMiddleware(ctx *context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("appContext", ctx)
			return next(c)
		}
	}
}

// ConfigMiddleware configura os middlewares a serem utilizados pelo servidor.
func ConfigMiddleware(e *echo.Echo, ctx *context.Context) {
	ctx.Logger.Info("Configurando middlewares")

	corsConfig := middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     ctx.Config.Origins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PATCH,
			echo.DELETE,
			echo.HEAD,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}

	// Redirecionamento para HTTPS
	e.Pre(middleware.HTTPSRedirect())

	e.Use(
		// Implementar app.AppWrapper
		ContextMiddleware(ctx),
		// Middleware para capturar requisições e respostas
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			// Filtrar campos para não aparecer nos logs
			filtered := make([][]byte, 2)
			bodies := [][]byte{reqBody, resBody}
			for j, body := range bodies {
				var i interface{}
				if err := json.Unmarshal(body, &i); err != nil {
					filtered[j] = []byte("")
				}

				if m, ok := i.(map[string]interface{}); ok {
					delete(m, "password")
					delete(m, "content")
					delete(m, "blob")
				} else {
					filtered[j] = []byte("")
				}

				if payload, err := json.Marshal(i); err == nil {
					filtered[j] = payload
				} else {
					filtered[j] = []byte("")
				}
			}

			handlers.LogHTTPDetails(
				c,
				zapcore.InfoLevel,
				"HTTP request-response",
				zap.Int("status", c.Response().Status),
				zap.String("request_body", string(filtered[0])),
				zap.String("response_body", string(filtered[1])),
			)
		}),
		// Limitações de requisições IP/segundo
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(rate.Limit(20)),
		),
		// XSSProtection
		middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "DENY",
			HSTSMaxAge:            31536000, // 1 ano
			HSTSPreloadEnabled:    true,
			ContentSecurityPolicy: "default-src 'self'; script-src 'self'; style-src 'self'",
		}),
		middleware.Recover(),
		// CORS
		middleware.CORSWithConfig(corsConfig),
		// Sistema de arquivos estáticos
		middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.Dir("frontend/dist"),
		}),
	)
}

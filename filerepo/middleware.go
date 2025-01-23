package filerepo

import (
	"agros_arquivos_patrocinadoras/filerepo/handlers"
	"agros_arquivos_patrocinadoras/filerepo/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"net/http"
)

// ContextMiddleware é o middleware para implementar services.AppWrapper como
// contexto padrão a ser usado por echo.Echo.
func ContextMiddleware(ctx *services.AppWrapper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("appContext", ctx)
			return next(c)
		}
	}
}

// ConfigMiddleware configura os middlewares a serem utilizados pelo servidor.
func ConfigMiddleware(e *echo.Echo, ctx *services.AppWrapper) {
	ctx.Logger.Info("Configurando middlewares")

	corsConfig := middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     ctx.Config.Origins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PATCH,
			echo.PUT,
			echo.DELETE,
			echo.HEAD,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}

	// Redirecionamento para HTTPS
	e.Pre(middleware.HTTPSRedirect())

	e.Use(
		// Middleware para capturar requisições e respostas
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			handlers.LogHTTPDetails(
				c,
				ctx.Logger,
				zapcore.InfoLevel,
				"HTTP request-response",
				zap.Int("status", c.Response().Status),
				zap.String("request_body", string(reqBody)),
				zap.String("response_body", string(resBody)),
			)
		}),
		// Implementar app.AppWrapper
		ContextMiddleware(ctx),
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
		// JWT
		//echojwt.WithConfig(echojwt.Config{
		//	NewClaimsFunc: func(c echo.AppWrapper) jwt.Claims {
		//		return new(auth.CustomClaims)
		//	},
		//	SigningKey: []byte(ctx.Config.JwtSecret),
		//}),
		// Sistema de arquivos estáticos
		middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.Dir("frontend/dist"),
		}),
	)
}

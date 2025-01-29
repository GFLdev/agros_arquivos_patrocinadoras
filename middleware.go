package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/auth"
	"agros_arquivos_patrocinadoras/pkg/handlers"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func AuthenticationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Ignorar rota [POST] /login
			// FIXME: Retirar 'c.Path() == "/user"'
			if (c.Path() == "/login" || c.Path() == "/user") && c.Request().Method == echo.POST {
				return next(c)
			}

			// Obter o token JWT do cabeçalho
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				return echo.NewHTTPError(http.StatusUnauthorized, handlers.UnauthorizedMessage)
			}
			tokenString := authHeader[7:]

			// Obter o contexto da aplicação para acessar a chave secreta
			ctx := context.GetContext(c)

			// Validar e analisar o token
			token, err := jwt.ParseWithClaims(
				tokenString,
				&auth.CustomClaims{},
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("método de assinatura inesperado: %w", token.Header["alg"])
					}
					return []byte(ctx.Config.JwtSecret), nil
				},
			)
			if err != nil || !token.Valid {
				ctx.Logger.Error("Token inválido.", zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, handlers.UnauthorizedMessage)
			}

			// Obter claims
			claims, ok := token.Claims.(*auth.CustomClaims)
			if !ok || claims == nil || claims.Id == uuid.Nil {
				return echo.NewHTTPError(http.StatusUnauthorized, handlers.UnauthorizedMessage)
			}

			// Adicionar informações do usuário ao contexto
			user := echo.Map{
				"userId": claims.Id,
				"name":   claims.Name,
				"admin":  claims.Admin,
			}
			c.Set("user", user)

			// Continuar para próximo handler
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
		// Implementar app.AppWrapper
		ContextMiddleware(ctx),
		// Middleware para capturar requisições e respostas
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			handlers.LogHTTPDetails(
				c,
				zapcore.InfoLevel,
				"HTTP request-response",
				zap.Int("status", c.Response().Status),
				zap.String("request_body", string(reqBody)),
				zap.String("response_body", string(resBody)),
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
		// Autenticação
		AuthenticationMiddleware(),
		// Sistema de arquivos estáticos
		middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: http.Dir("frontend/dist"),
		}),
	)
}

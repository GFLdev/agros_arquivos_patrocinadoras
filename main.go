package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"strconv"
)

// Serve inicializa o servidor echo.
func Serve(e *echo.Echo, ctx *AppContext) {
	var err error
	if ctx.Config.Environment == "production" {
		ctx.Logger.Info(fmt.Sprintf("Iniciando servidor de produção na porta %d", ctx.Config.Port))

		err = e.StartAutoTLS(":" + strconv.Itoa(ctx.Config.Port))
	} else if ctx.Config.Environment == "development" {
		ctx.Logger.Info(fmt.Sprintf("Iniciando servidor de desenvolvimento na porta %d", ctx.Config.Port))

		err = e.StartTLS(
			":"+strconv.Itoa(ctx.Config.Port),
			ctx.Config.CertFile,
			ctx.Config.KeyFile,
		)
	} else {
		err = fmt.Errorf("erro na configuração de ambiente")
	}

	if err != nil {
		ctx.Logger.Fatal("Erro na execução do servidor",
			zap.Error(err),
		)
	}
}

func init() {
	// Criação da pasta logs, caso não exista
	err := os.MkdirAll("logs", 0777)
	if err != nil {
		panic(err)
	}

	// Criação da pasta com os arquivos, caso não exista
	err = os.MkdirAll("db", 0777)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Contexto da aplicação
	ctx := &AppContext{
		Logger: CreateLogger(),
	}
	ctx.Logger.Info("Iniciando aplicação")

	// Configurações
	config := LoadConfig(ctx)
	ctx.Config = config

	// Servidor Echo
	e := echo.New()
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	ConfigMiddleware(e, ctx)

	// Rotas
	ConfigRoutes(e, ctx)

	// Inicializar servidor
	Serve(e, ctx)
}

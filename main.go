package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const (
	FsRoot = "data"
)

// Serve inicializa o servidor echo.
func Serve(e *echo.Echo, ctx *context.Context) {
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

func handleSIGINT(c chan os.Signal, logr *zap.Logger) {
	for sig := range c {
		if sig == syscall.SIGINT {
			logr.Warn("SIGINT recebido")

			var i string
			fmt.Print("Deseja finalizar a aplicação? [S/n] (default: n) ")
			_, err := fmt.Scanln(&i)
			if err == nil && strings.ToUpper(i) == "S" {
				logr.Info("Finalizando a aplicação")
				os.Exit(0)
			} else {
				logr.Info("SIGINT interrompido")
			}
		}
	}
}

func firstSetup(logr *zap.Logger) {
	logr.Info("Configurando aplicação pela primeira vez")

}

func init() {
	// Criação da pasta logs, caso não exista
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Criação da pasta com os arquivos, caso não exista
	err = os.MkdirAll(FsRoot, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Logger
	logr := logger.CreateLogger()

	logr.Info("Iniciando aplicação")

	// Configurações
	cfg := config.LoadConfig(logr)

	// Sistema de arquivos
	filesystem := &fs.FileSystem{
		Root: FsRoot,
	}

	// Banco de dados
	dataBase := db.GetSqlDB(&cfg.Database, logr)

	// Contexto da aplicação
	ctx := &context.Context{
		Logger:     logr,
		Config:     cfg,
		FileSystem: filesystem,
		DB:         dataBase,
	}

	// Servidor Echo
	e := echo.New()
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	ConfigMiddleware(e, ctx)

	// Rotas
	ConfigRoutes(e, ctx)

	// Handler para SIGINT (^C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go handleSIGINT(c, logr)

	// Inicializar servidor
	Serve(e, ctx)
}

package main

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/services/config"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"agros_arquivos_patrocinadoras/filerepo/services/logger"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Serve inicializa o servidor echo.
func Serve(e *echo.Echo, ctx *services.AppWrapper) {
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

	// Criação da pasta com os arquivos, caso não exista
	err := os.MkdirAll("repo", os.ModePerm)
	if err != nil {
		logr.Fatal("Erro ao criar diretório do repositório", zap.Error(err))
	}

	// Novo repositório de arquivos
	newRepo := fs.FS{
		Users:     map[uuid.UUID]fs.User{},
		UpdatedAt: time.Now().Unix(),
	}

	// Escrita do arquivo de rastreamento dos arquivos
	err = fs.StructToFile[fs.FS]("repo/track.json", &newRepo, logr)
	if err != nil {
		logr.Fatal("Erro na escrita de repo/track.json", zap.Error(err))
	}
}

func init() {
	// Criação da pasta logs, caso não exista
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Logger
	logr := logger.CreateLogger()

	// Caso não exista o arquivo de rastreamento dos arquivos, execute a
	// primeira configuração da aplicação
	_, err = os.Stat("repo/track.json")
	if errors.Is(err, os.ErrNotExist) {
		firstSetup(logr)
	}
}

func main() {
	// Logger
	logr := logger.CreateLogger()

	logr.Info("Iniciando aplicação")

	// Configurações
	cfg := config.LoadConfig(logr)

	// Sistema de arquivos
	appFs, err := fs.GetFS(logr)
	if err != nil {
		logr.Fatal("Erro ao obter repositório de arquivos",
			zap.Error(err),
		)
	}

	// Banco de dados
	dataBase := db.GetSqlDB(cfg.Database, logr)

	// Contexto da aplicação
	ctx := services.NewAppWrapper(logr, cfg, appFs, dataBase)

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

package main

import (
	"agros_arquivos_patrocinadoras/config"
	"agros_arquivos_patrocinadoras/db"
	"agros_arquivos_patrocinadoras/handlers"
	"agros_arquivos_patrocinadoras/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"time"
)

// Serve inicializa o servidor echo.
func Serve(e *echo.Echo, ctx *handlers.AppContext) {
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

func firstSetup(logr *zap.Logger) {
	logr.Info("Configurando aplicação pela primeira vez")

	// Criação da pasta com os arquivos, caso não exista
	err := os.MkdirAll("repo", os.ModePerm)
	if err != nil {
		logr.Fatal("Erro ao criar diretório do repositório", zap.Error(err))
	}

	// Novo repositório de arquivos
	newRepo := db.Repo{
		Users:     map[uuid.UUID]db.User{},
		UpdatedAt: time.Now().Unix(),
	}

	// Escrita do arquivo de rastreamento dos arquivos
	err = db.StructToFile[db.Repo]("repo/track.json", &newRepo, logr)
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
	// Contexto da aplicação
	ctx := &handlers.AppContext{
		Logger: logger.CreateLogger(),
	}
	ctx.Logger.Info("Iniciando aplicação")

	// Configurações
	cfg := config.LoadConfig(ctx)
	ctx.Config = cfg

	// Repositório
	repo, err := db.GetFileRepo(ctx.Logger)
	if err != nil {
		ctx.Logger.Fatal("Erro ao obter repositório de arquivos",
			zap.Error(err),
		)
	}
	ctx.Repo = struct {
		*db.Repo
		*sync.Mutex
	}{repo, new(sync.Mutex)}

	// Servidor Echo
	e := echo.New()
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	handlers.ConfigMiddleware(e, ctx)

	// Rotas
	handlers.ConfigRoutes(e, ctx)

	// Inicializar servidor
	Serve(e, ctx)
}

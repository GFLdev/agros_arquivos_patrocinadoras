package main

import (
	"agros_arquivos_patrocinadoras/clogger"
	config2 "agros_arquivos_patrocinadoras/config"
	"agros_arquivos_patrocinadoras/db"
	"agros_arquivos_patrocinadoras/handlers"
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

func firstSetup(logger *zap.Logger) {
	logger.Info("Configurando aplicação pela primeira vez")

	// Criação da pasta com os arquivos, caso não exista
	err := os.MkdirAll("repo", os.ModePerm)
	if err != nil {
		logger.Fatal("Erro ao criar diretório do repositório", zap.Error(err))
	}

	// Novo repositório de arquivos
	newRepo := db.Repo{
		Users:     map[uuid.UUID]db.User{},
		UpdatedAt: time.Now().Unix(),
	}

	// Escrita do arquivo de rastreamento dos arquivos
	err = db.StructToFile[db.Repo]("repo/track.json", &newRepo, logger)
	if err != nil {
		logger.Fatal("Erro na escrita de repo/track.json", zap.Error(err))
	}
}

func init() {
	// Criação da pasta logs, caso não exista
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Logger
	logger := clogger.CreateLogger()

	// Caso não exista o arquivo de rastreamento dos arquivos, execute a
	// primeira configuração da aplicação
	_, err = os.Stat("repo/track.json")
	if errors.Is(err, os.ErrNotExist) {
		firstSetup(logger)
	}
}

func main() {
	// Contexto da aplicação
	ctx := &handlers.AppContext{
		Logger: clogger.CreateLogger(),
	}
	ctx.Logger.Info("Iniciando aplicação")

	// Configurações
	config := config2.LoadConfig(ctx)
	ctx.Config = config

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

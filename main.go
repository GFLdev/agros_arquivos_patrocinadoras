// Package main implementa a aplicação principal responsável por inicializar e
// configurar o servidor HTTP, gerenciar sinais do sistema e configurar
// dependências como banco de dados, sistema de arquivos e logger.
package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// Serve é responsável por iniciar o servidor HTTP da aplicação.
//
// Parâmetros:
//   - e: instância do servidor Echo para configurar e iniciar
//     o serviço HTTP/HTTPS.
//   - ctx: contexto da aplicação contendo configurações,
//     logger e dependências essenciais.
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

// handleSIGINT é responsável por lidar com o sinal SIGINT (Ctrl+C) recebido
// durante a execução da aplicação, permitindo que o usuário decida se deseja
// finalizar ou continuar a execução.
//
// Parâmetros:
//   - c: canal para capturar os sinais enviados ao processo.
//   - logr: instância do logger usada para registrar informações e avisos
//     durante o manuseio do sinal.
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
			}
			logr.Info("SIGINT interrompido")
		}
	}
}

func init() {
	// Criação da pasta logs, caso não exista
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		panic(err)
	}
}

func main() {
	// Logger
	logr := logger.CreateLogger()

	// Handler para SIGINT (^C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go handleSIGINT(c, logr)

	logr.Info("Iniciando aplicação")

	// Configurações
	cfg := config.LoadConfig(logr)

	// Banco de dados
	dataBase := db.GetSqlDB(&cfg.Database, logr)
	defer func(dataBase *sql.DB) {
		err := dataBase.Close()
		if err != nil {
			logr.Error("Erro ao fechar banco de dados", zap.Error(err))
		}
	}(dataBase)

	// Contexto da aplicação
	ctx := &context.Context{
		Logger: logr,
		Config: cfg,
		DB:     dataBase,
	}

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

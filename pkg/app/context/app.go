package context

import (
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"database/sql"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Context contém as informações e recursos necessários para o processamento
// de uma solicitação, incluindo o logger, configuração, sistema de arquivos e
// banco de dados.
type Context struct {
	// Logger é o logger usado para registrar informações, erros e eventos
	// durante o processamento da solicitação.
	Logger *zap.Logger
	// Config contém as configurações necessárias para o aplicativo,
	// como credenciais de banco de dados, parâmetros de ambiente, etc.
	Config *config.Config
	// FileSystem é o sistema de arquivos que gerencia operações de leitura
	// e escrita de arquivos durante o processamento.
	FileSystem *fs.FileSystem
	// DB é a conexão com o banco de dados, usada para executar operações de
	// consulta e modificação de dados.
	DB *sql.DB
}

// GetContext recupera o app.Context associado à solicitação a partir do
// objeto echo.Context.
//
// Retorna o contexto associado à requisição atual.
func GetContext(c echo.Context) *Context {
	return c.Get("appContext").(*Context)
}

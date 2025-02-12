// Package context fornece uma estrutura para gerenciar as informações
// necessárias ao processamento de solicitações, integrando recursos como
// logger, configuração, sistema de arquivos e conexão com banco de dados.
package context

import (
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"database/sql"
	"github.com/google/uuid"
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
	// DB é a conexão com o banco de dados, usada para executar operações de
	// consulta e modificação de dados.
	DB      *sql.DB
	AdminId uuid.UUID
}

// GetContext retorna o contexto da aplicação a partir do contexto da
// solicitação do Echo.
//
// Parâmetros:
//   - c: contexto da solicitação do Echo, que armazena os dados da requisição.
//
// Retorno:
//   - *Context: ponteiro para o contexto da aplicação associado à solicitação.
func GetContext(c echo.Context) *Context {
	return c.Get("appContext").(*Context)
}

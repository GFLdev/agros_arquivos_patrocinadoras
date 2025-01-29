// Package db fornece funcionalidades para conexão com o banco de dados Oracle
// utilizando a biblioteca database/sql.
package db

import (
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"go.uber.org/zap"
)

// GetSqlDB estabelece uma conexão com o banco de dados Oracle.
//
// Parâmetros:
//   - dbParams: ponteiro para uma struct config.Database contendo os parâmetros
//     de conexão ao banco de dados (username, senha, servidor, porta e serviço).
//   - logr: instância do logger zap usada para registrar mensagens de log.
//
// Retorno:
//   - *sql.DB: instância do banco de dados conectada.
//
// Observação: O programa é encerrado com uma mensagem de erro fatal caso falhe
// ao estabelecer ou testar a conexão.
func GetSqlDB(dbParams *config.Database, logr *zap.Logger) *sql.DB {
	// String de conexão
	connString := fmt.Sprintf(
		"oracle://%s:%s@%s:%s/%s",
		dbParams.Username,
		dbParams.Password,
		dbParams.Server,
		dbParams.Port,
		dbParams.Service,
	)

	// Conectar ao banco
	logr.Info("Conectando ao banco de dados")
	db, err := sql.Open("oracle", connString)
	if err != nil {
		logr.Fatal("Erro ao conectar ao banco de dados", zap.Error(err))
	}

	// Testar conexão
	logr.Info("Testando conexão ao banco de dados")
	if err = db.Ping(); err != nil {
		logr.Fatal("Erro ao testar a conexão ao banco de dados", zap.Error(err))
	}
	logr.Info("Banco de dados conectado com sucesso")
	return db
}

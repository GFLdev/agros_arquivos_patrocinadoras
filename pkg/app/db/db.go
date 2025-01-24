package db

import (
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"go.uber.org/zap"
)

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
		logr.Fatal("Erro ao conectar ao banco de dados",
			zap.Error(err),
		)
	}

	// Testar conexão
	logr.Info("Testando conexão ao banco de dados")
	err = db.Ping()
	if err != nil {
		logr.Fatal("Erro ao testar a conexão ao banco de dados",
			zap.Error(err),
		)
	}

	logr.Info("Banco de dados conectado com sucesso")
	return db
}

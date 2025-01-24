package db

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// rollbackCreate realiza o rollback de uma transação e, caso necessário, exclui
// um arquivo ou diretório. Essa função é chamada quando ocorre um erro durante
// a criação de uma entidade no banco de dados ou no sistema de arquivos.
//
// Ela tenta reverter as ações de banco de dados e do sistema de arquivos (caso
// tenha ocorrido alguma criação), garantindo que o sistema volte ao estado
// anterior em caso de falhas.
//
// Parâmetros:
//
// - tx: transação aberta no banco de dados que deve ser revertida.
//
// - errDb: erro ocorrido na operação de banco de dados que exige o rollback.
//
// - errFs: erro ocorrido na operação de sistema de arquivos que exige a
// exclusão de arquivos ou diretórios.
//
// - path: caminho do arquivo ou diretório que deve ser removido caso tenha sido
// criado com falha. Caso não haja caminho, essa parte da operação não será
// executada.
//
// Retorno:
//
// - Não retorna valores. Realiza ações de rollback e limpeza, mas não gera erro
// diretamente.
func rollbackCreate(tx *sql.Tx, errDb, errFs error, path string) {
	if tx != nil && (errDb != nil || errFs != nil) {
		// Logs
		logr := logger.CreateLogger()

		// Tentativas de rollback no banco
		err := tx.Rollback()
		if err != nil {
			logr.Error(
				"Tentativa de rollback falhou",
				zap.Error(err),
			)
		}

		// Tentativa de exclusão do arquivo/diretório, caso tenha sido criado
		if path != "" {
			err = fs.DeleteFile(path)
			if err != nil {
				logr.Error(
					"Tentativa de limpeza falhou",
					zap.Error(err),
				)
			}
		}
	}
}

func CreateUser(ctx *context.Context, p UserCreation) error {
	var errDb, errFs error
	var path string

	// Iniciar uma transação
	tx, errDb := ctx.DB.Begin()
	if errDb != nil {
		return fmt.Errorf("não foi possível transação: %v", errDb)
	}

	// Agendar rollback em caso de erro
	defer rollbackCreate(tx, errDb, errFs, path)

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	userId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível transação: %v", errDb)
	}

	// Insert query
	schema := ctx.Config.Database.Schema
	insert := fmt.Sprintf(
		`INSERT INTO %s.%s
  		(%s, %s, %s, %s)
		VALUES ($1, $2, $3, $4)`,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.UserId,
		schema.UserTable.Columns.Name,
		schema.UserTable.Columns.Password,
		schema.UserTable.Columns.UpdatedAt,
	)

	// Criação
	_, errDb = tx.Exec(insert, userId.String(), p.Name, p.Password, ts)
	if errDb != nil {
		return fmt.Errorf("não foi possível criar usuário: %v", errDb)
	}

	// Transação no sistema de arquivos
	path = fmt.Sprintf("%s/%s", ctx.FileSystem.Root, userId.String())
	errFs = ctx.FileSystem.CreateEntity(path, nil, fs.User)
	if errFs != nil {
		return fmt.Errorf("não foi possível diretório do usuário: %v", errFs)
	}

	// Confirmar a transação no banco
	errDb = tx.Commit()
	if errDb != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", errDb)
	}

	return nil
}

func CreateCategory(
	db *sql.DB,
	fs *fs.FileSystem,
	userId uuid.UUID,
	categId uuid.UUID,
	name string,
) error {
	var errDb, errFs error
	var path string

	// Iniciar uma transação
	tx, errDb := db.Begin()
	if errDb != nil {
		return fmt.Errorf("não foi possível transação: %v", errDb)
	}

	// Agendar rollback em caso de erro
	defer rollbackCreate(tx, errDb, errFs, path)

	// Insert
	insert := `INSERT INTO patrocinadoras_categoria
  		(uuid_categ, uuid_patroc, nome_categ, ts_modificado)
		VALUES ($1, $2, $3, $4)`

	// Criação
	_, errDb = tx.Exec(insert, categId.String(), userId.String(), name, time.Now().Unix())
	if errDb != nil {
		return fmt.Errorf("não foi possível criar categoria: %v", errDb)
	}

	// Transação no sistema de arquivos
	path, errFs = fs.CreateCategory(userId, categId)
	if errFs != nil {
		return fmt.Errorf("não foi possível diretório da categoria: %v", errFs)
	}

	// Confirmar a transação no banco
	errDb = tx.Commit()
	if errDb != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", errDb)
	}

	return nil
}

func CreateFile(
	db *sql.DB,
	fs *fs.FileSystem,
	fileId uuid.UUID,
	userId uuid.UUID,
	categId uuid.UUID,
	name string,
	extension string,
	mimeType string,
	content *[]byte,
) error {
	var errDb, errFs error
	var path string

	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("não foi possível transação: %v", err)
	}

	// Agendar rollback em caso de erro
	defer rollbackCreate(tx, errDb, errFs, path)

	// Insert
	insert := `INSERT INTO patrocinadoras_arquivo
  		(, , , , , )
		VALUES ($1, $2, $3, $4, $5, $6)`
	ts := time.Now().Unix()

	// Criação
	_, err = tx.Exec(
		insert,
		fileId,
		categId,
		name,
		extension,
		mimeType,
		ts,
	)
	if err != nil {
		return fmt.Errorf("não foi possível criar arquivo: %v", err)
	}

	// Transação no sistema de arquivos
	path, errFs = fs.CreateFile(userId, categId, fileId, extension, context)
	if errFs != nil {
		return fmt.Errorf("não foi possível diretório da categoria: %v", errFs)
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", err)
	}

	return nil
}

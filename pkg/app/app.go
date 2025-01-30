// Package app fornece funcionalidades essenciais para a aplicação, incluindo
// operações relacionadas a gerenciamento de usuários, interações com o banco
// de dados, manipulação de sistema de arquivos e controle de transações. Ele
// centraliza os componentes principais utilizados em várias partes da
// aplicação, promovendo reutilização de código e consistência nas operações.
package app

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/types/db"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CheckUsername(ctx *context.Context, name string) (bool, error) {
	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s
		FROM %s.%s
		WHERE %s = :name`,
		schema.UserTable.Columns.UserId,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.Name,
	)

	// Obtenção da linha
	var userId uuid.UUID
	row := ctx.DB.QueryRow(query, sql.Named("name", name))
	err := row.Scan(&userId)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, fmt.Errorf("nome de usuário já existente")
}

// GetCredentials realiza a busca de informações de login de um usuário no banco
// de dados com base em seu nome.
//
// Parâmetros:
//   - ctx: contexto da aplicação contendo recursos necessários como o banco
//     de dados e configurações.
//   - name: string representando o nome do usuário a ser procurado.
//
// Retorno:
//   - LoginCompare: estrutura contendo o ID do usuário e o hash da senha.
//   - error: erro que pode ocorrer durante a busca, como falhas de query ou
//     ausência de informações.
func GetCredentials(ctx *context.Context, p UserParams) (uuid.UUID, error) {
	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s, %s
		FROM %s.%s
		WHERE %s = :name`,
		schema.UserTable.Columns.UserId,
		schema.UserTable.Columns.Password,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.Name,
	)

	// Obtenção da linha
	rows, err := ctx.DB.Query(query, sql.Named("name", p.Name))
	if err != nil {
		return uuid.Nil, fmt.Errorf("não foi possível obter os usuários")
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		l := LoginCompare{}
		err = rows.Scan(&l.UserId, &l.Hash)
		if err != nil {
			continue
		}
		err = bcrypt.CompareHashAndPassword([]byte(l.Hash), []byte(p.Password))
		if err == nil {
			return l.UserId, nil
		}
	}

	return uuid.Nil, fmt.Errorf("não autenticado")
}

// rollbackCreate realiza um rollback em caso de falha durante a criação de
// uma entidade no banco de dados e/ou no sistema de arquivos.
//
// Parâmetros:
//   - ctx: contexto da aplicação que contém recursos compartilhados como
//     banco de dados, logger e sistema de arquivos.
//   - rbData: estrutura CreateRollbackData contendo informações da transação,
//     caminho no sistema de arquivos e captura de erros durante o rollback.
func rollbackCreate(ctx *context.Context, rbData CreateRollbackData) {
	if rbData.Tx != nil && (rbData.DB != nil || rbData.FS != nil) {
		// Tentativas de rollback no banco
		if rbData.DB != nil {
			if err := rbData.Tx.Rollback(); err != nil {
				ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(err))
			}
		}

		// Tentativa de exclusão do arquivo/diretório, caso tenha sido criado
		if rbData.Path != "" && rbData.FS != nil {
			err := ctx.FileSystem.DeleteEntity(rbData.Path)
			if err != nil {
				ctx.Logger.Error("Tentativa de limpeza falhou", zap.Error(err))
			}
		}
	}
}

// CreateUser cria um registro de um novo usuário no banco de dados e
// configura um diretório correspondente no sistema de arquivos.
//
// Parâmetros:
//   - ctx: ponteiro para o contexto da aplicação, contendo recursos como o
//     banco de dados, sistema de arquivos, e informações de configuração.
//   - p: estrutura UserParams com os dados necessários para criar o usuário
//     (ex: nome, senha).
//
// Retorno:
//   - uuid.UUID: UUID do usuário criado.
//   - error: retorna qualquer erro que possa ocorrer durante o processo de
//     criação, como falhas na transação ou erros na interação com o sistema
//     de arquivos.
func CreateUser(ctx *context.Context, p UserParams) (uuid.UUID, error) {
	rbErrors := &RollbackErrors{}

	// Checar nome de usuário
	ok, err := CheckUsername(ctx, p.Name)
	if !ok {
		ctx.Logger.Error("Erro no nome de usuário.", zap.Error(err))
		return uuid.Nil, err
	}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	userId, err := uuid.NewUUID()
	if err != nil {
		ctx.Logger.Error("Erro ao criar UUID.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar UUID")
	}
	path := filepath.Join(ctx.FileSystem.Root, userId.String())

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	rbData := CreateRollbackData{
		Tx:             tx,
		Path:           path,
		RollbackErrors: rbErrors,
	}
	defer rollbackCreate(ctx, rbData)

	// Insert query
	schema := &ctx.Config.Database.Schema
	insert := fmt.Sprintf(
		`INSERT INTO %s.%s
  		(%s, %s, %s, %s)
		VALUES (:user_id, :name, :password, :updated_at)`,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.UserId,
		schema.UserTable.Columns.Name,
		schema.UserTable.Columns.Password,
		schema.UserTable.Columns.UpdatedAt,
	)

	// Criptografar senha
	hash, err := HashPassword(ctx, p.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("não foi possível criptografar senha")
	}

	// Criação
	_, rbErrors.DB = tx.Exec(
		insert,
		sql.Named("user_id", userId.String()),
		sql.Named("name", p.Name),
		sql.Named("password", hash),
		sql.Named("updated_at", ts),
	)
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar usuário.", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível criar usuário")
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.CreateEntity(path, nil, fs.User)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao criar diretório.", zap.Error(rbErrors.FS))
		return uuid.Nil, rbErrors.FS
	}

	// Confirmar a transação no banco
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível confirmar transação")
	}
	return userId, nil
}

// CreateCategory cria um registro de uma nova categoria no banco de dados e
// configura um diretório correspondente no sistema de arquivos.
//
// Parâmetros:
//   - ctx: ponteiro para o contexto da aplicação, contendo o banco de dados,
//     sistema de arquivos e configurações.
//   - p: estrutura CategParams com os dados necessários para criar a categoria
//     (ex: nome, identificador do usuário).
//
// Retorno:
//   - uuid.UUID: UUID da categoria criada.
//   - error: retorna um erro caso alguma etapa do processo falhe, como falhas
//     na transação ou erros na criação do diretório.
func CreateCategory(ctx *context.Context, p CategParams) (uuid.UUID, error) {
	rbErrors := &RollbackErrors{}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	categId, err := uuid.NewUUID()
	if err != nil {
		ctx.Logger.Error("Erro ao criar UUID.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar UUID")
	}
	path := filepath.Join(ctx.FileSystem.Root, p.UserId.String(), categId.String())

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	rbData := CreateRollbackData{
		Tx:             tx,
		Path:           path,
		RollbackErrors: rbErrors,
	}
	defer rollbackCreate(ctx, rbData)

	// Insert query
	schema := &ctx.Config.Database.Schema
	insert := fmt.Sprintf(
		`INSERT INTO %s.%s
  		(%s, %s, %s, %s)
		VALUES (:categ_id, :user_id, :name, :updated_at)`,
		schema.Name,
		schema.CategTable.Name,
		schema.CategTable.Columns.CategId,
		schema.CategTable.Columns.UserId,
		schema.CategTable.Columns.Name,
		schema.CategTable.Columns.UpdatedAt,
	)

	// Criação
	_, rbErrors.DB = tx.Exec(
		insert,
		sql.Named("categ_id", categId.String()),
		sql.Named("user_id", p.UserId.String()),
		sql.Named("name", p.Name),
		sql.Named("updated_at", ts),
	)
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar categoria.", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível criar categoria")
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.CreateEntity(path, nil, fs.Category)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao criar diretório.", zap.Error(rbErrors.FS))
		return uuid.Nil, rbErrors.FS
	}

	// Confirmar a transação no banco
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível confirmar transação")
	}
	return categId, nil
}

// CreateFile cria um registro de um novo arquivo no banco de dados e
// configura um diretório correspondente no sistema de arquivos.
//
// Parâmetros:
//   - ctx: ponteiro para o contexto da aplicação, que contém o banco de dados,
//     sistema de arquivos e as configurações da aplicação.
//   - p: estrutura FileParams que contém os dados necessários para a criação do
//     arquivo (ex: nome, extensão, conteúdo, identificador da categoria e do
//     usuário).
//
// Retorno:
//   - uuid.UUID: UUID do arquivo criado.
//   - error: retorna um erro caso alguma etapa do processo falhe, seja no banco
//     de dados ou no sistema de arquivos.
func CreateFile(ctx *context.Context, p FileParams) (uuid.UUID, error) {
	rbErrors := &RollbackErrors{}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	fileId, err := uuid.NewUUID()
	if err != nil {
		ctx.Logger.Error("Erro ao criar UUID.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar UUID")
	}
	path := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
		fileId.String()+p.Extension,
	)

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	rbData := CreateRollbackData{
		Tx:             tx,
		Path:           path,
		RollbackErrors: rbErrors,
	}
	defer rollbackCreate(ctx, rbData)

	// Insert query
	schema := &ctx.Config.Database.Schema
	insert := fmt.Sprintf(
		`INSERT INTO %s.%s
  		(%s, %s, %s, %s, %s, %s)
		VALUES (:file_id, :categ_id, :name, :extension, :mimetype, :updated_at)`,
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.FileId,
		schema.FileTable.Columns.CategId,
		schema.FileTable.Columns.Name,
		schema.FileTable.Columns.Extension,
		schema.FileTable.Columns.Mimetype,
		schema.FileTable.Columns.UpdatedAt,
	)

	// Criação
	_, rbErrors.DB = tx.Exec(
		insert,
		sql.Named("file_id", fileId.String()),
		sql.Named("categ_id", p.CategId.String()),
		sql.Named("name", p.Name),
		sql.Named("extension", p.Extension),
		sql.Named("mimetype", p.Mimetype),
		sql.Named("updated_at", ts),
	)
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar arquivo.", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível criar arquivo")
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.CreateEntity(path, p.Content, fs.File)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao criar arquivo em disco.", zap.Error(rbErrors.FS))
		return uuid.Nil, rbErrors.FS
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return uuid.Nil, fmt.Errorf("não foi possível confirmar transação")
	}
	return fileId, nil
}

// closeRows fecha as linhas abertas de uma consulta SQL para liberar os
// recursos no banco de dados.
//
// Parâmetros:
//   - ctx: o contexto da aplicação, contendo o Logger para registrar
//     mensagens de advertência.
//   - rows: ponteiro para o resultado da consulta SQL, que será fechado.
func closeRows(ctx *context.Context, rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		ctx.Logger.Warn("Erro ao fechar linhas da query", zap.Error(err))
	}
}

// QueryAllUsers recupera todos os usuários armazenados no banco de dados.
//
// Parâmetros:
//   - ctx: o contexto da aplicação, contendo a configuração do banco de dados
//     e o zap.Logger para registrar logs de advertência em caso de erro.
//
// Retorno:
//   - []db.UserModel: uma lista de usuários contendo os campos UserId, Name e
//     UpdatedAt.
//   - error: um erro é retornado caso a query ou o processamento dos resultados
//     falhe.
func QueryAllUsers(ctx *context.Context) ([]db.UserModel, error) {
	var users []db.UserModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s FROM %s.%s`,
		schema.UserTable.Columns.UserId,
		schema.UserTable.Columns.Name,
		schema.UserTable.Columns.UpdatedAt,
		schema.Name,
		schema.UserTable.Name,
	)

	// Obtenção das linhas
	rows, err := ctx.DB.Query(query)
	if err != nil {
		return users, fmt.Errorf("não foi possível obter os usuários")
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		u := db.UserModel{Password: ""}
		err = rows.Scan(&u.UserId, &u.Name, &u.UpdatedAt)
		if err != nil {
			ctx.Logger.Error("Erro ao obter usuário.", zap.Error(err))
			return users, fmt.Errorf("não foi possível obter todos os usuários")
		}
		users = append(users, u)
	}
	return users, nil
}

// QueryAllCategories recupera todas as categorias associadas a um usuário
// específico do banco de dados.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo a configuração do banco de dados
//     e o zap.Logger.
//   - userId: identificador único do usuário cujas categorias devem ser
//     recuperadas.
//
// Retorno:
//   - []db.CategModel: uma lista de categorias contendo os campos CategId,
//     UserId, Name e UpdatedAt.
//   - error: um erro é retornado caso a consulta ou o processamento dos
//     resultados falhe.
func QueryAllCategories(ctx *context.Context, userId uuid.UUID) ([]db.CategModel, error) {
	var categs []db.CategModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s,%s
		FROM %s.%s
		WHERE %s = :user_id`,
		schema.CategTable.Columns.CategId,
		schema.CategTable.Columns.UserId,
		schema.CategTable.Columns.Name,
		schema.CategTable.Columns.UpdatedAt,
		schema.Name,
		schema.CategTable.Name,
		schema.CategTable.Columns.UserId,
	)

	// Obtenção das linhas
	rows, err := ctx.DB.Query(query, sql.Named("user_id", userId.String()))
	if err != nil {
		return categs, fmt.Errorf("não foi possível obter as categorias")
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		var c db.CategModel
		err = rows.Scan(&c.CategId, &c.UserId, &c.Name, &c.UpdatedAt)
		if err != nil {
			ctx.Logger.Error("Erro ao obter categoria.", zap.Error(err))
			return categs, fmt.Errorf("não foi possível obter todas as categorias")
		}
		categs = append(categs, c)
	}
	return categs, nil
}

// QueryAllFiles recupera todos os arquivos associados a uma categoria
// específica do banco de dados.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo a configuração do banco de dados
//     e o zap.Logger.
//   - categId: identificador único da categoria cujos arquivos devem ser
//     recuperados.
//
// Retorno:
//   - []db.FileModel: uma lista de arquivos contendo os campos FileId,
//     CategId, Name, Extension, Mimetype e UpdatedAt.
//   - error: um erro é retornado caso a consulta ou o processamento dos
//     resultados falhe.
func QueryAllFiles(ctx *context.Context, categId uuid.UUID) ([]db.FileModel, error) {
	var files []db.FileModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s,%s,%s,%s
		FROM %s.%s
		WHERE %s = :categ_id`,
		schema.FileTable.Columns.FileId,
		schema.FileTable.Columns.CategId,
		schema.FileTable.Columns.Name,
		schema.FileTable.Columns.Extension,
		schema.FileTable.Columns.Mimetype,
		schema.FileTable.Columns.UpdatedAt,
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.CategId,
	)

	// Obtenção das linhas
	rows, err := ctx.DB.Query(query, sql.Named("categ_id", categId.String()))
	if err != nil {
		return files, fmt.Errorf("não foi possível obter os arquivos")
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		var f db.FileModel
		err = rows.Scan(
			&f.FileId,
			&f.CategId,
			&f.Name,
			&f.Extension,
			&f.Mimetype,
			&f.UpdatedAt,
		)
		if err != nil {
			ctx.Logger.Error("Erro ao obter arquivo.", zap.Error(err))
			return files, fmt.Errorf("não foi possível obter todas os arquivos")
		}
		files = append(files, f)
	}

	return files, nil
}

// QueryUserById realiza uma consulta ao banco de dados para buscar um
// usuário pelo seu ID.
//
// Parâmetros:
//   - ctx: contexto da aplicação contendo configurações e acesso ao banco
//     de dados.
//   - userId: o uuid.UUID do usuário a ser buscado.
//
// Retorno:
//   - db.UserModel: estrutura contendo os dados do usuário encontrado.
//   - error: retorna um erro caso a execução da consulta ou o
//     processamento do resultado falhe.
func QueryUserById(ctx *context.Context, userId uuid.UUID) (db.UserModel, error) {
	var user db.UserModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s
		FROM %s.%s
		WHERE %s = :user_id`,
		schema.UserTable.Columns.UserId,
		schema.UserTable.Columns.Name,
		schema.UserTable.Columns.UpdatedAt,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.UserId,
	)

	// Obtenção da linha
	row := ctx.DB.QueryRow(query, sql.Named("user_id", userId.String()))
	err := row.Scan(&user.UserId, &user.Name, &user.UpdatedAt)
	if err != nil {
		return user, fmt.Errorf("não foi possível obter usuário")
	}
	user.Password = ""
	return user, nil
}

// QueryCategoryById realiza uma consulta ao banco de dados para buscar uma
// categoria pelo seu ID.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo informações de configuração e
//     acesso ao banco de dados.
//   - categId: o uuid.UUID da categoria a ser buscada.
//
// Retorno:
//   - db.CategModel: estrutura contendo os dados da categoria encontrada.
//   - error: retorna um erro caso ocorra falha na execução da consulta ou no
//     processamento do resultado.
func QueryCategoryById(ctx *context.Context, categId uuid.UUID) (db.CategModel, error) {
	var categ db.CategModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s,%s
		FROM %s.%s
		WHERE %s = :categ_id`,
		schema.CategTable.Columns.CategId,
		schema.CategTable.Columns.UserId,
		schema.CategTable.Columns.Name,
		schema.CategTable.Columns.UpdatedAt,
		schema.Name,
		schema.CategTable.Name,
		schema.CategTable.Columns.CategId,
	)

	// Obtenção da linha
	row := ctx.DB.QueryRow(query, sql.Named("categ_id", categId.String()))
	err := row.Scan(&categ.CategId, &categ.UserId, &categ.Name, &categ.UpdatedAt)
	if err != nil {
		return categ, fmt.Errorf("não foi possível obter categoria")
	}
	return categ, nil
}

// QueryFileById realiza uma consulta ao banco de dados para buscar um arquivo
// pelo seu ID.
//
// Parâmetros:
//   - ctx: contexto da aplicação contendo informações de configuração e acesso
//     ao banco de dados.
//   - fileId: o uuid.UUID do arquivo a ser buscado.
//
// Retorno:
//   - db.FileModel: estrutura contendo os dados do arquivo encontrado.
//   - error: retorna um erro caso ocorra falha na execução da consulta ou no
//     processamento do resultado.
func QueryFileById(ctx *context.Context, fileId uuid.UUID) (db.FileModel, error) {
	var file db.FileModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s,%s,%s,%s
		FROM %s.%s
		WHERE %s = :file_id`,
		schema.FileTable.Columns.FileId,
		schema.FileTable.Columns.CategId,
		schema.FileTable.Columns.Name,
		schema.FileTable.Columns.Extension,
		schema.FileTable.Columns.Mimetype,
		schema.FileTable.Columns.UpdatedAt,
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.FileId,
	)

	// Obtenção da linha
	row := ctx.DB.QueryRow(query, sql.Named("file_id", fileId.String()))
	err := row.Scan(
		&file.FileId,
		&file.CategId,
		&file.Name,
		&file.Extension,
		&file.Mimetype,
		&file.UpdatedAt,
	)
	if err != nil {
		return file, fmt.Errorf("não foi possível obter arquivo")
	}
	return file, nil
}

// rollbackUpdate tenta reverter uma atualização em caso de falha, restaurando o
// estado anterior no banco de dados e no sistema de arquivos, caso necessário.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo configurações, conexão ao banco de
//     dados e utilitários de log.
//   - rbData: estrutura UpdateRollbackData com informações da transação a ser
//     revertida, caminhos de arquivos (se aplicável) e erros de rollback.
func rollbackUpdate(ctx *context.Context, rbData UpdateRollbackData) {
	if rbData.Tx != nil && (rbData.DB != nil || rbData.FS != nil) {
		// Tentativas de rollback no banco
		if rbData.DB != nil {
			if err := rbData.Tx.Rollback(); err != nil {
				ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(err))
			}
		}

		// Tentativa de mover arquivo para o caminho original, caso tenha sido
		// movido
		if rbData.OldPath != "" && rbData.NewPath != "" && rbData.FS != nil {
			err := ctx.FileSystem.UpdateEntity(rbData.NewPath, rbData.OldPath)
			if err != nil {
				ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(err))
			}
		}
	}
}

// UpdateUser atualiza os dados de um usuário no banco e, opcionalmente, no
// sistema de arquivos, garantindo consistência por meio de transações. Apenas os
// parâmetros inseridos são alterados.
//
// Parâmetros:
//   - ctx: contexto da aplicação contendo configurações, conexão com o banco de
//     dados e utilitários para log.
//   - p: estrutura UserUpdate com os detalhes do usuário a ser atualizado,
//     incluindo IDs, nomes (atual e novo), senha, entre outros.
//
// Retorno:
//   - error: retorna um erro caso ocorra falha em qualquer etapa do processo,
//     seja na inicialização ou confirmação da transação, ou na atualização no
//     banco de dados.
func UpdateUser(ctx *context.Context, p UserUpdate) error {
	rbErrors := &RollbackErrors{}

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Geração do Timestamp
	ts := time.Now().Unix()

	// Agendar rollback em caso de erro
	rbData := UpdateRollbackData{
		Tx:             tx,
		OldPath:        "",
		NewPath:        "",
		RollbackErrors: rbErrors,
	}
	defer rollbackUpdate(ctx, rbData)

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []any
	var set []string
	if p.Name != p.OldName {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.UserTable.Columns.Name+" = :name")
	}
	if p.Password != "" {
		// Criptografar senha
		hash, err := HashPassword(ctx, p.Password)
		if err != nil {
			return fmt.Errorf("não foi possível criptografar senha")
		}
		args = append(args, sql.Named("password", hash))
		set = append(set, schema.UserTable.Columns.Password+" = :password")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.UserTable.Columns.UpdatedAt+" = :updated_at")

	// Update query
	update := fmt.Sprintf(`UPDATE %s.%s
				SET %s
				WHERE %s = :user_id`,
		schema.Name,
		schema.UserTable.Name,
		strings.Join(set, ","),
		schema.UserTable.Columns.UserId,
	)
	args = append(args, sql.Named("user_id", p.UserId.String()))

	// Atualização
	var res sql.Result
	res, rbErrors.DB = tx.Exec(update, args...)
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao atualizar usuário.", zap.Error(rbErrors.DB))
		return fmt.Errorf("não foi possível atualizar usuário")
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao atualizar usuário.", zap.Error(rbErrors.DB))
		return rbErrors.DB
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

// UpdateCategory atualiza os dados de uma categoria no banco e, opcionalmente, no
// sistema de arquivos, garantindo consistência por meio de transações. Apenas os
// parâmetros inseridos são alterados.
//
// Parâmetros:
//   - ctx: contexto da aplicação contendo configurações, conexão com o banco de
//     dados, utilitários de log e acesso ao sistema de arquivos.
//   - p: estrutura CategUpdate com informações da categoria a ser atualizada,
//     incluindo IDs, nomes (atual e novo) e usuários.
//
// Retorno:
//   - error: retorna um erro caso ocorra falha em qualquer etapa, seja na
//     inicialização ou confirmação da transação, execução da alteração no banco
//     de dados ou atualização no sistema de arquivos.
func UpdateCategory(ctx *context.Context, p CategUpdate) error {
	rbErrors := &RollbackErrors{}

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Geração do Timestamp e caminhos
	ts := time.Now().Unix()
	oldPath := filepath.Join(
		ctx.FileSystem.Root,
		p.OldUserId.String(),
		p.CategId.String(),
	)
	newPath := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
	)

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []any
	var set []string
	if oldPath != newPath {
		if !ctx.FileSystem.EntityExists(filepath.Dir(newPath)) {
			rbErrors.DB = fmt.Errorf(
				"não pôde atualizar categoria: diretório do usuário %s não existe",
				filepath.Dir(newPath),
			)
			return rbErrors.DB
		}
		args = append(args, sql.Named("user_id", p.UserId.String()))
		set = append(set, schema.CategTable.Columns.UserId+" = :user_id")
	}
	if p.Name != p.OldName {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.CategTable.Columns.Name+" = :name")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.CategTable.Columns.UpdatedAt+" = :updated_at")

	// Agendar rollback em caso de erro
	rbData := UpdateRollbackData{
		Tx:             tx,
		OldPath:        oldPath,
		NewPath:        newPath,
		RollbackErrors: rbErrors,
	}
	defer rollbackUpdate(ctx, rbData)

	// Update query
	update := fmt.Sprintf(`UPDATE %s.%s
				SET %s
				WHERE %s = :categ_id`,
		schema.Name,
		schema.CategTable.Name,
		strings.Join(set, ","),
		schema.CategTable.Columns.CategId,
	)
	args = append(args, sql.Named("categ_id", p.CategId.String()))

	// Atualização
	var res sql.Result
	res, rbErrors.DB = tx.Exec(update, args...)
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao atualizar categoria.", zap.Error(rbErrors.DB))
		return fmt.Errorf("não foi possível atualizar categoria")
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao atualizar categoria.", zap.Error(rbErrors.DB))
		return rbErrors.DB
	}

	// Transação no sistema de arquivos
	if oldPath != newPath {
		rbErrors.FS = ctx.FileSystem.UpdateEntity(oldPath, newPath)
		if rbErrors.FS != nil {
			return rbErrors.FS
		}
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

// UpdateFile atualiza os dados de um arquivo no banco e, opcionalmente, no
// sistema de arquivos, garantindo consistência por meio de transações. Apenas os
// parâmetros inseridos são alterados.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo configurações, conexão com o banco,
//     sistema de arquivos e utilitários.
//   - p: estrutura FileUpdate com os dados do arquivo a ser atualizado,
//     incluindo ID, categoria, nome, extensão, mime type, entre outros.
//
// Retorno:
//   - error: retorna erro se ocorrer falha durante a transação ou ao aplicar
//     mudanças no banco, ou no sistema de arquivos.
func UpdateFile(ctx *context.Context, p FileUpdate) error {
	rbErrors := &RollbackErrors{}

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Geração do Timestamp e caminhos
	ts := time.Now().Unix()
	oldPath := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.OldCategId.String(),
		p.FileId.String()+p.OldExtension,
	)
	newPath := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
		p.FileId.String(),
	)

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []any
	var set []string
	if oldPath != newPath {
		if !ctx.FileSystem.EntityExists(filepath.Dir(newPath)) {
			rbErrors.DB = fmt.Errorf(
				"não pôde atualizar arquivo: diretório %s a categoria não existe",
				filepath.Dir(newPath),
			)
			return rbErrors.DB
		}
		args = append(args, sql.Named("categ_id", p.CategId.String()))
		set = append(set, schema.FileTable.Columns.CategId+" = :categ_id")
	}
	if p.Name != p.OldName {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.FileTable.Columns.Name+" = :name")
	}
	if p.Extension != p.OldExtension {
		args = append(args, sql.Named("extension", p.Extension))
		set = append(set, schema.FileTable.Columns.Extension+" = :extension")
		newPath = newPath + p.Extension
	} else {
		newPath = newPath + p.OldExtension
	}
	if p.Mimetype != p.OldMimetype {
		args = append(args, sql.Named("mimetype", p.Mimetype))
		set = append(set, schema.FileTable.Columns.Mimetype+" = :mimetype")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.CategTable.Columns.UpdatedAt+" = :updated_at")

	// Agendar rollback em caso de erro
	rbData := UpdateRollbackData{
		Tx:             tx,
		OldPath:        oldPath,
		NewPath:        newPath,
		RollbackErrors: rbErrors,
	}
	defer rollbackUpdate(ctx, rbData)

	// Update query
	update := fmt.Sprintf(`UPDATE %s.%s
				SET %s
				WHERE %s = :file_id`,
		schema.Name,
		schema.FileTable.Name,
		strings.Join(set, ","),
		schema.FileTable.Columns.FileId,
	)
	args = append(args, sql.Named("file_id", p.FileId.String()))

	// Atualização
	var res sql.Result
	res, rbErrors.DB = tx.Exec(update, args...)
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao atualizar arquivo.", zap.Error(rbErrors.DB))
		return fmt.Errorf("não foi possível atualizar arquivo")
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao atualizar arquivo.", zap.Error(rbErrors.DB))
		return rbErrors.DB
	}

	// Transação no sistema de arquivos
	if p.Content == nil && oldPath != newPath {
		rbErrors.FS = ctx.FileSystem.UpdateEntity(oldPath, newPath)
		if rbErrors.FS != nil {
			return rbErrors.FS
		}
	} else if p.Content != nil {
		rbErrors.FS = ctx.FileSystem.CreateEntity(newPath, p.Content, fs.File)
		if rbErrors.FS != nil {
			return rbErrors.FS
		}
		defer func(ctx *context.Context, rbData UpdateRollbackData) {
			if rbData.FS != nil {
				// Tentativa de exclusão do arquivo criado, caso tenha erros
				err := ctx.FileSystem.DeleteEntity(rbData.NewPath)
				if err != nil {
					ctx.Logger.Error("Tentativa de limpeza falhou", zap.Error(err))
				}
			} else {
				// Tentativa de exclusão do arquivo antigo, quando não tiver erros
				err := ctx.FileSystem.DeleteEntity(rbData.OldPath)
				if err != nil {
					ctx.Logger.Error("Tentativa de limpeza falhou", zap.Error(err))
				}
			}
		}(ctx, rbData)
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

// rollbackDelete realiza o rollback de uma transação de exclusão no banco de
// dados e no sistema de arquivos.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo informações de configuração, banco
//     de dados e sistema de arquivos.
//   - rbData: estrutura DeleteRollbackData com os dados necessários para
//     executar o rollback, como a transação, caminho do arquivo e erros.
func rollbackDelete(ctx *context.Context, rbData DeleteRollbackData) {
	if rbData.Tx != nil && (rbData.DB != nil || rbData.FS != nil) {
		// Tentativas de rollback no banco
		if rbData.DB != nil {
			if err := rbData.Tx.Rollback(); err != nil {
				ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(err))
			}
		}

		// Tentativa de incluir o arquivo backup
		if rbData.FS != nil {
			err := ctx.FileSystem.CreateEntity(
				rbData.Path,
				rbData.Content,
				rbData.Type,
			)
			if err != nil {
				ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(err))
			}
		}
	}
}

// DeleteUser realiza a exclusão de um usuário no banco de dados e no sistema de
// arquivos.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo informações de configuração, banco de
//     dados e sistema de arquivos.
//   - p: estrutura UserDelete que contém as informações do usuário a ser excluído.
//
// Retorno:
//   - error: retorna um erro se ocorrer alguma falha durante o processo de
//     exclusão.
func DeleteUser(ctx *context.Context, p UserDelete) error {
	rbErrors := &RollbackErrors{}

	// Caminho
	path := filepath.Join(ctx.FileSystem.Root, p.UserId.String())

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Agendar rollback em caso de erro
	rbData := DeleteRollbackData{
		Tx:             tx,
		Path:           path,
		Type:           fs.User,
		Content:        nil,
		RollbackErrors: rbErrors,
	}
	defer rollbackDelete(ctx, rbData)

	// Delete query
	schema := &ctx.Config.Database.Schema
	del := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE %s = :user_id",
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.UserId,
	)

	// Exclusão
	var res sql.Result
	res, rbErrors.DB = tx.Exec(del, sql.Named("user_id", p.UserId.String()))
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao excluir usuário.", zap.Error(rbErrors.DB))
		return fmt.Errorf("não foi possível excluir usuário")
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao excluir usuário.", zap.Error(rbErrors.DB))
		return rbErrors.DB
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao confirmar transação")
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.DeleteEntity(path)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao excluir diretório.", zap.Error(rbErrors.DB))
		return rbErrors.FS
	}
	return nil
}

// DeleteCategory realiza a exclusão de uma categoria no banco de dados e no
// sistema de arquivos.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo configurações, banco de dados e
//     sistema de arquivos.
//   - p: estrutura CategDelete com os dados necessários para a exclusão.
//
// Retorno:
//   - error: retorna um erro se ocorrer falha durante o processo de exclusão.
func DeleteCategory(ctx *context.Context, p CategDelete) error {
	rbErrors := &RollbackErrors{}

	// Caminho
	path := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
	)

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Agendar rollback em caso de erro
	rbData := DeleteRollbackData{
		Tx:             tx,
		Path:           path,
		Type:           fs.Category,
		Content:        nil,
		RollbackErrors: rbErrors,
	}
	defer rollbackDelete(ctx, rbData)

	// Delete query
	schema := &ctx.Config.Database.Schema
	del := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE %s = :categ_id",
		schema.Name,
		schema.CategTable.Name,
		schema.CategTable.Columns.CategId,
	)

	// Exclusão
	var res sql.Result
	res, rbErrors.DB = tx.Exec(del, sql.Named("categ_id", p.CategId.String()))
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao excluir categoria.", zap.Error(rbErrors.DB))
		return fmt.Errorf("não foi possível excluir categoria")
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao excluir categoria.", zap.Error(rbErrors.DB))
		return rbErrors.DB
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao confirmar transação")
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.DeleteEntity(path)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao excluir diretório.", zap.Error(rbErrors.DB))
		return rbErrors.FS
	}
	return nil
}

// DeleteFile realiza a exclusão de um arquivo no banco de dados e no sistema
// de arquivos.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo configurações, banco de dados e
//     sistema de arquivos.
//   - p: estrutura FileDelete com os dados necessários para a exclusão.
//
// Retorno:
//   - error: retorna um erro se ocorrer falha durante o processo de exclusão.
func DeleteFile(ctx *context.Context, p FileDelete) error {
	rbErrors := &RollbackErrors{}

	// Caminho
	path := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
		p.FileId.String()+p.Extension,
	)

	// Obter conteúdo do arquivo para backup
	// Abrir arquivo
	file, err := os.Open(path)
	if err != nil {
		ctx.Logger.Error("Erro ao abrir arquivo.", zap.Error(err))
		rbErrors.FS = fmt.Errorf("erro ao abrir %s", path)
		return rbErrors.FS
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			ctx.Logger.Error("Erro ao fechar arquivo para leitura", zap.Error(err))
		}
	}(file)
	// Leitura
	var backupContent []byte
	backupContent, rbErrors.FS = io.ReadAll(file)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao ler arquivo.", zap.Error(rbErrors.FS))
		rbErrors.FS = fmt.Errorf("erro ao ler %s para backup", path)
		return rbErrors.FS
	}

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Agendar rollback em caso de erro
	rbData := DeleteRollbackData{
		Tx:             tx,
		Path:           path,
		Type:           fs.File,
		Content:        &backupContent,
		RollbackErrors: rbErrors,
	}
	defer rollbackDelete(ctx, rbData)

	// Delete query
	schema := &ctx.Config.Database.Schema
	del := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE %s = :file_id",
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.FileId,
	)

	// Exclusão
	var res sql.Result
	res, rbErrors.DB = tx.Exec(del, sql.Named("file_id", p.FileId.String()))
	if rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao excluir arquivo.", zap.Error(rbErrors.DB))
		return fmt.Errorf("não foi possível excluir arquivo")
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao excluir arquivo.", zap.Error(rbErrors.DB))
		return rbErrors.DB
	}

	// Confirmar a transação
	if rbErrors.DB = tx.Commit(); rbErrors.DB != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(rbErrors.DB))
		return fmt.Errorf("erro ao confirmar transação")
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.DeleteEntity(path)
	if rbErrors.FS != nil {
		ctx.Logger.Error("Erro ao excluir arquivo em disco.", zap.Error(rbErrors.DB))
		return rbErrors.FS
	}

	return nil
}

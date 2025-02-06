// Package app fornece funcionalidades essenciais para a aplicação, incluindo
// operações relacionadas a gerenciamento de usuários, interações com o banco
// de dados, manipulação de sistema de arquivos e controle de transações. Ele
// centraliza os componentes principais utilizados em várias partes da
// aplicação, promovendo reutilização de código e consistência nas operações.
package app

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/types/db"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	goora "github.com/sijms/go-ora/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
		ctx.Logger.Error("Erro ao buscar usuário", zap.Error(err))
		return false, fmt.Errorf("não foi possível procurar usuário")
	}
	return false, fmt.Errorf("nome de usuário já existente")
}

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
		return uuid.Nil, fmt.Errorf("usuário não encontrado")
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

func rollback(ctx *context.Context, tx *sql.Tx, err *error) {
	if tx != nil && *err != nil {
		// Tentativa de rollback no banco
		if *err = tx.Rollback(); *err != nil {
			ctx.Logger.Error("Tentativa de rollback falhou", zap.Error(*err))
		}
	}
}

func CreateUser(ctx *context.Context, p UserParams) (uuid.UUID, error) {
	var err error

	// Checar nome de usuário
	ok, err := CheckUsername(ctx, p.Name)
	if !ok {
		return uuid.Nil, err
	}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	userId, err := uuid.NewUUID()
	if err != nil {
		ctx.Logger.Error("Erro ao criar UUID.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar UUID")
	}

	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

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
		return uuid.Nil, err
	}

	// Criação
	_, err = tx.Exec(
		insert,
		sql.Named("user_id", userId.String()),
		sql.Named("name", p.Name),
		sql.Named("password", hash),
		sql.Named("updated_at", ts),
	)
	if err != nil {
		ctx.Logger.Error("Erro ao criar usuário.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar usuário")
	}

	// Confirmar a transação no banco
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível confirmar transação")
	}
	return userId, nil
}

func CreateCategory(ctx *context.Context, p CategParams) (uuid.UUID, error) {
	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	categId, err := uuid.NewUUID()
	if err != nil {
		ctx.Logger.Error("Erro ao criar UUID.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar UUID")
	}

	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

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
	_, err = tx.Exec(
		insert,
		sql.Named("categ_id", categId.String()),
		sql.Named("user_id", p.UserId.String()),
		sql.Named("name", p.Name),
		sql.Named("updated_at", ts),
	)
	if err != nil {
		ctx.Logger.Error("Erro ao criar categoria.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar categoria")
	}

	// Confirmar a transação no banco
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível confirmar transação")
	}
	return categId, nil
}

func CreateFile(ctx *context.Context, p FileParams) (uuid.UUID, error) {
	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	fileId, err := uuid.NewUUID()
	if err != nil {
		ctx.Logger.Error("Erro ao criar UUID.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar UUID")
	}

	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar transação")
	}

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

	// Criação do BLOB
	directLob := goora.Blob{Data: *p.Content}

	// Insert query
	schema := &ctx.Config.Database.Schema
	insert := fmt.Sprintf(
		`INSERT INTO %s.%s
  		(%s, %s, %s, %s, %s, %s, %s)
		VALUES (:file_id, :categ_id, :name, :extension, :mimetype, :blob, :updated_at)`,
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.FileId,
		schema.FileTable.Columns.CategId,
		schema.FileTable.Columns.Name,
		schema.FileTable.Columns.Extension,
		schema.FileTable.Columns.Mimetype,
		schema.FileTable.Columns.Blob,
		schema.FileTable.Columns.UpdatedAt,
	)

	// Criação
	_, err = tx.Exec(
		insert,
		sql.Named("file_id", fileId.String()),
		sql.Named("categ_id", p.CategId.String()),
		sql.Named("name", p.Name),
		sql.Named("extension", p.Extension),
		sql.Named("mimetype", p.Mimetype),
		sql.Named("blob", directLob),
		sql.Named("updated_at", ts),
	)
	if err != nil {
		ctx.Logger.Error("Erro ao criar arquivo.", zap.Error(err))
		return uuid.Nil, fmt.Errorf("não foi possível criar arquivo")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
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
		f.Blob = nil
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
		`SELECT %s,%s,%s,%s,%s,%s,%s
		FROM %s.%s
		WHERE %s = :file_id`,
		schema.FileTable.Columns.FileId,
		schema.FileTable.Columns.CategId,
		schema.FileTable.Columns.Name,
		schema.FileTable.Columns.Extension,
		schema.FileTable.Columns.Mimetype,
		schema.FileTable.Columns.Blob,
		schema.FileTable.Columns.UpdatedAt,
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.FileId,
	)

	// Obtenção da linha
	var directLob goora.Blob
	row := ctx.DB.QueryRow(query, sql.Named("file_id", fileId.String()))
	err := row.Scan(
		&file.FileId,
		&file.CategId,
		&file.Name,
		&file.Extension,
		&file.Mimetype,
		&directLob,
		&file.UpdatedAt,
	)
	if err != nil {
		return file, fmt.Errorf("não foi possível obter arquivo")
	}
	file.Blob = directLob.Data
	return file, nil
}

func UpdateUser(ctx *context.Context, userId uuid.UUID, p UserParams) error {
	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Geração do Timestamp
	ts := time.Now().Unix()

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []any
	var set []string
	if p.Name != "" {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.UserTable.Columns.Name+" = :name")
	}
	if p.Password != "" {
		// Criptografar senha
		var hash string
		hash, err = HashPassword(ctx, p.Password)
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
	args = append(args, sql.Named("user_id", userId.String()))

	// Atualização
	res, err := tx.Exec(update, args...)
	if err != nil {
		ctx.Logger.Error("Erro ao atualizar usuário.", zap.Error(err))
		return fmt.Errorf("não foi possível atualizar usuário")
	} else if n, _ := res.RowsAffected(); n > 1 {
		err = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao atualizar usuário.", zap.Error(err))
		return fmt.Errorf("não foi possível atualizar usuário")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

func UpdateCategory(ctx *context.Context, categId uuid.UUID, p CategParams) error {
	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Geração do Timestamp
	ts := time.Now().Unix()

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []any
	var set []string
	if p.UserId != uuid.Nil {
		args = append(args, sql.Named("user_id", p.UserId.String()))
		set = append(set, schema.CategTable.Columns.UserId+" = :user_id")
	}
	if p.Name != "" {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.CategTable.Columns.Name+" = :name")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.CategTable.Columns.UpdatedAt+" = :updated_at")

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

	// Update query
	update := fmt.Sprintf(`UPDATE %s.%s
				SET %s
				WHERE %s = :categ_id`,
		schema.Name,
		schema.CategTable.Name,
		strings.Join(set, ","),
		schema.CategTable.Columns.CategId,
	)
	args = append(args, sql.Named("categ_id", categId.String()))

	// Atualização
	res, err := tx.Exec(update, args...)
	if err != nil {
		ctx.Logger.Error("Erro ao atualizar categoria.", zap.Error(err))
		return fmt.Errorf("não foi possível atualizar categoria")
	} else if n, _ := res.RowsAffected(); n > 1 {
		err = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao atualizar categoria.", zap.Error(err))
		return fmt.Errorf("não foi possível atualizar categoria")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

func UpdateFile(ctx *context.Context, fileId uuid.UUID, p FileParams) error {
	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Geração do Timestamp
	ts := time.Now().Unix()

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []any
	var set []string
	if p.CategId != uuid.Nil {
		args = append(args, sql.Named("categ_id", p.CategId.String()))
		set = append(set, schema.FileTable.Columns.CategId+" = :categ_id")
	}
	if p.Name != "" {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.FileTable.Columns.Name+" = :name")
	}
	if p.Extension != "" && p.Extension != "." {
		args = append(args, sql.Named("extension", p.Extension))
		set = append(set, schema.FileTable.Columns.Extension+" = :extension")
	}
	if p.Mimetype != "" {
		args = append(args, sql.Named("mimetype", p.Mimetype))
		set = append(set, schema.FileTable.Columns.Mimetype+" = :mimetype")
	}
	if p.Content != nil && len(*p.Content) > 0 {
		directLob := goora.Blob{Data: *p.Content}
		args = append(args, sql.Named("blob", directLob))
		set = append(set, schema.FileTable.Columns.Blob+" = :blob")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.CategTable.Columns.UpdatedAt+" = :updated_at")

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

	// Update query
	update := fmt.Sprintf(`UPDATE %s.%s
				SET %s
				WHERE %s = :file_id`,
		schema.Name,
		schema.FileTable.Name,
		strings.Join(set, ","),
		schema.FileTable.Columns.FileId,
	)
	args = append(args, sql.Named("file_id", fileId.String()))

	// Atualização
	res, err := tx.Exec(update, args...)
	if err != nil {
		ctx.Logger.Error("Erro ao atualizar arquivo.", zap.Error(err))
		return fmt.Errorf("não foi possível atualizar arquivo")
	} else if n, _ := res.RowsAffected(); n > 1 {
		err = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao atualizar arquivo.", zap.Error(err))
		return fmt.Errorf("não foi possível atualizar arquivo")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

func DeleteUser(ctx *context.Context, userId uuid.UUID) error {
	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

	// Delete query
	schema := &ctx.Config.Database.Schema
	del := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE %s = :user_id",
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.UserId,
	)

	// Exclusão
	res, err := tx.Exec(del, sql.Named("user_id", userId.String()))
	if err != nil {
		ctx.Logger.Error("Erro ao excluir usuário.", zap.Error(err))
		return fmt.Errorf("não foi possível excluir usuário")
	} else if n, _ := res.RowsAffected(); n > 1 {
		err = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao excluir usuário.", zap.Error(err))
		return fmt.Errorf("não foi possível excluir usuário")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

func DeleteCategory(ctx *context.Context, categId uuid.UUID) error {
	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

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
	res, err = tx.Exec(del, sql.Named("categ_id", categId.String()))
	if err != nil {
		ctx.Logger.Error("Erro ao excluir categoria.", zap.Error(err))
		return fmt.Errorf("não foi possível excluir categoria")
	} else if n, _ := res.RowsAffected(); n > 1 {
		err = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao excluir categoria.", zap.Error(err))
		return fmt.Errorf("não foi possível excluir categoria")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

func DeleteFile(ctx *context.Context, fileId uuid.UUID) error {
	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		ctx.Logger.Error("Erro ao criar transação de banco.", zap.Error(err))
		return fmt.Errorf("erro ao iniciar transação")
	}

	// Agendar rollback em caso de erro
	defer rollback(ctx, tx, &err)

	// Delete query
	schema := &ctx.Config.Database.Schema
	del := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE %s = :file_id",
		schema.Name,
		schema.FileTable.Name,
		schema.FileTable.Columns.FileId,
	)

	// Exclusão
	res, err := tx.Exec(del, sql.Named("file_id", fileId.String()))
	if err != nil {
		ctx.Logger.Error("Erro ao excluir arquivo.", zap.Error(err))
		return fmt.Errorf("não foi possível excluir arquivo")
	} else if n, _ := res.RowsAffected(); n > 1 {
		err = fmt.Errorf("mais de uma linha afetada")
		ctx.Logger.Error("Erro ao excluir arquivo.", zap.Error(err))
		return fmt.Errorf("não foi possível excluir arquivo")
	}

	// Confirmar a transação
	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("Erro ao efetivar transação (COMMIT).", zap.Error(err))
		return fmt.Errorf("erro ao confirmar transação")
	}
	return nil
}

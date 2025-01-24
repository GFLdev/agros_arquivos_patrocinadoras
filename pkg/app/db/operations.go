package db

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/types/db"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

// rollbackCreate executa as ações de rollback para desfazer alterações parciais
// realizadas durante uma operação.
//
// Este método tenta reverter a transação do banco de dados e excluir arquivos
// ou diretórios criados no sistema de arquivos em caso de falha em qualquer
// etapa do processo.
//
// Parâmetros:
//
// - ctx: contexto da aplicação contendo o zap.Logger para registro de erros e o
// sistema de arquivos para remover o arquivo/diretório em path.
//
// - tx: ponteiro para a transação do banco de dados que será revertida.
//
// - errDb: ponteiro para o erro relacionado ao banco de dados, indicando falha
// na operação.
//
// - errFs: ponteiro para o erro relacionado ao sistema de arquivos, indicando
// falha na criação ou exclusão de arquivos/diretórios.
//
// - path: ponteiro para o caminho do arquivo ou diretório que deve ser excluído
// no rollback.
func rollbackCreate(
	ctx *context.Context,
	tx *sql.Tx,
	errDb, errFs *error,
	path *string,
) {
	if tx != nil && (*errDb != nil || *errFs != nil) {
		// Tentativas de rollback no banco
		if *errDb != nil {
			err := tx.Rollback()
			if err != nil {
				ctx.Logger.Error(
					"Tentativa de rollback falhou",
					zap.Error(err),
				)
			}
		}

		// Tentativa de exclusão do arquivo/diretório, caso tenha sido criado
		if *path != "" && *errFs != nil {
			err := ctx.FileSystem.DeleteEntity(*path)
			if err != nil {
				ctx.Logger.Error(
					"Tentativa de limpeza falhou",
					zap.Error(err),
				)
			}
		}
	}
}

// CreateUser cria um novo usuário no banco de dados e um diretório associado no
// sistema de arquivos.
//
// Este método executa uma transação que envolve a criação de um registro no
// banco de dados e a criação de um diretório correspondente para o usuário no
// sistema de arquivos. Em caso de falha, um rollback é executado para reverter
// quaisquer alterações realizadas.
//
// Parâmetros:
//
// - ctx: ponteiro para o contexto da aplicação, que contém o banco de dados,
// configurações e sistema de arquivos.
//
// - p: estrutura UserCreation contendo as informações do usuário a ser criado.
//
// Retorno:
//
// - error: retorna um erro caso ocorra qualquer falha durante o processo de
// criação, incluindo erros de transação no banco de dados ou no sistema de
// arquivos.
func CreateUser(ctx *context.Context, p *UserCreation) error {
	var errDb, errFs error
	var path string

	// Iniciar uma transação
	tx, errDb := ctx.DB.Begin()
	if errDb != nil {
		return fmt.Errorf("não foi possível transação: %v", errDb)
	}

	// Agendar rollback em caso de erro
	defer rollbackCreate(ctx, tx, &errDb, &errFs, &path)

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	userId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível criar UUID: %v", errDb)
	}

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

	// Criação
	_, errDb = tx.Exec(insert, userId.String(), p.Name, p.Password, ts)
	if errDb != nil {
		return fmt.Errorf("não foi possível criar usuário: %v", errDb)
	}

	// Transação no sistema de arquivos
	path = filepath.Join(ctx.FileSystem.Root, userId.String())
	errFs = ctx.FileSystem.CreateEntity(path, nil, fs.User)
	if errFs != nil {
		return fmt.Errorf("não foi possível criar diretório %s: %v", path, errFs)
	}

	// Confirmar a transação no banco
	errDb = tx.Commit()
	if errDb != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", errDb)
	}

	return nil
}

// CreateCategory cria uma nova categoria associada a um usuário no banco de
// dados e um diretório correspondente no sistema de arquivos.
//
// Este método executa uma transação que envolve a criação de um registro no
// banco de dados para a categoria e a criação de um diretório no sistema de
// arquivos. Em caso de falha, um rollback é executado para reverter quaisquer
// alterações realizadas.
//
// Parâmetros:
//
// - ctx: ponteiro para o contexto da aplicação, que contém o banco de dados,
// configurações e sistema de arquivos.
//
// - p: estrutura CategCreation contendo as informações da categoria a ser
// criada.
//
// Retorno:
//
// - error: retorna um erro caso ocorra qualquer falha durante o processo de
// criação, incluindo erros de transação no banco de dados ou no sistema de
// arquivos.
func CreateCategory(ctx *context.Context, p *CategCreation) error {
	var errDb, errFs error
	var path string

	// Iniciar uma transação
	tx, errDb := ctx.DB.Begin()
	if errDb != nil {
		return fmt.Errorf("não foi possível transação: %v", errDb)
	}

	// Agendar rollback em caso de erro
	defer rollbackCreate(ctx, tx, &errDb, &errFs, &path)

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	categId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível criar UUID: %v", errDb)
	}

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
	_, errDb = tx.Exec(insert, categId.String(), p.UserId.String(), p.Name, ts)
	if errDb != nil {
		return fmt.Errorf("não foi possível criar categoria: %v", errDb)
	}

	// Transação no sistema de arquivos
	path = filepath.Join(ctx.FileSystem.Root, p.UserId.String(), categId.String())
	errFs = ctx.FileSystem.CreateEntity(path, nil, fs.Category)
	if errFs != nil {
		return fmt.Errorf("não foi possível criar diretório %s: %v", path, errFs)
	}

	// Confirmar a transação no banco
	errDb = tx.Commit()
	if errDb != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", errDb)
	}

	return nil
}

// CreateFile cria um novo registro de arquivo no banco de dados e um arquivo
// correspondente no sistema de arquivos.
//
// O método executa uma transação que envolve a criação de um registro no banco
// de dados e a criação física de um arquivo. Em caso de erro, realiza um
// rollback para desfazer alterações parciais.
//
// Parâmetros:
//
// - ctx: ponteiro para o contexto da aplicação, que contém o banco de dados,
// configurações e sistema de arquivos.
//
// - p: estrutura FileCreation contendo as informações do arquivo a ser criado.
//
// Retorno:
//
// - error: retorna um erro caso ocorra falha durante o processo de criação no
// banco de dados ou no sistema de arquivos.
func CreateFile(ctx *context.Context, p *FileCreation) error {
	var errDb, errFs error
	var path string

	// Iniciar uma transação
	tx, err := ctx.DB.Begin()
	if err != nil {
		return fmt.Errorf("não foi possível transação: %v", err)
	}

	// Agendar rollback em caso de erro
	defer rollbackCreate(ctx, tx, &errDb, &errFs, &path)

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	fileId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível criar UUID: %v", errDb)
	}

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
	_, err = tx.Exec(
		insert,
		fileId.String(),
		p.CategId.String(),
		p.Name,
		p.Extension,
		p.Mimetype,
		ts,
	)
	if err != nil {
		return fmt.Errorf("não foi possível criar arquivo: %v", err)
	}

	// Transação no sistema de arquivos
	path = filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
		fileId.String()+p.Extension,
	)
	errFs = ctx.FileSystem.CreateEntity(path, nil, fs.File)
	if errFs != nil {
		return fmt.Errorf("não foi possível criar arquivo %s: %v", path, errFs)
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", err)
	}

	return nil
}

// closeRows garante o fechamento de um ponteiro de linhas de resultado de uma
// consulta SQL.
//
// Essa função é utilizada para evitar vazamento de recursos, garantindo que o
// objeto sql.Rows seja fechado corretamente após o uso. Caso ocorra um erro
// durante o fechamento, ele será registrado no zap.Logger.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, usado para registrar logs de aviso em caso de
// falha ao fechar as linhas.
//
// - rows: ponteiro para o objeto sql.Rows que será fechado.
func closeRows(ctx *context.Context, rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		ctx.Logger.Warn(
			"Erro ao fechar linhas da query",
			zap.Error(err),
		)
	}
}

// QueryAllUsers recupera todos os usuários do banco de dados.
//
// Esta função executa uma consulta SQL para buscar todos os registros da tabela
// de usuários e os retorna como uma lista de modelos db.UserModel. As senhas
// não são incluídas nos resultados por segurança.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo a configuração do banco de dados e o
// zap.Logger.
//
// Retorno:
//
// - []db.UserModel: uma lista de usuários recuperados do banco, contendo os
// campos UserId, Name e UpdatedAt.
//
// - error: um erro é retornado caso a consulta ou o processamento dos
// resultados falhe.
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
		return users, fmt.Errorf(
			"não foi possível realizar query de todos os usuários: %s",
			err,
		)
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		u := db.UserModel{Password: ""}
		err := rows.Scan(&u.UserId, &u.Name, &u.UpdatedAt)
		if err != nil {
			return users, fmt.Errorf(
				"não foi possível obter todos os usuários: %v",
				err,
			)
		}
		users = append(users, u)
	}

	return users, nil
}

// QueryAllCategories recupera todas as categorias associadas a um usuário
// específico no banco de dados.
//
// Esta função executa uma consulta SQL para buscar todas as categorias
// vinculadas ao userId fornecido e retorna os resultados como uma lista de
// modelos db.CategModel.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo a configuração do banco de dados e o
// zap.Logger.
//
// - userId: identificador único do usuário cujas categorias devem ser
// recuperadas.
//
// Retorno:
//
// - []db.CategModel: uma lista de categorias contendo os campos CategId,
// UserId, Name e UpdatedAt.
//
// - error: um erro é retornado caso a consulta ou o processamento dos
// resultados falhe.
func QueryAllCategories(
	ctx *context.Context,
	userId uuid.UUID,
) ([]db.CategModel, error) {
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
	rows, err := ctx.DB.Query(query, userId.String())
	if err != nil {
		return categs, fmt.Errorf(
			"não foi possível realizar query de todas as categorias: %s",
			err,
		)
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		var c db.CategModel
		err := rows.Scan(&c.CategId, &c.UserId, &c.Name, &c.UpdatedAt)
		if err != nil {
			return categs, fmt.Errorf(
				"não foi possível obter todas as categorias: %v",
				err,
			)
		}
		categs = append(categs, c)
	}

	return categs, nil
}

// QueryAllFiles recupera todos os arquivos associados a uma categoria
// específica no banco de dados.
//
// Esta função executa uma consulta SQL para buscar todos os arquivos vinculados
// ao categId fornecido e retorna os resultados como uma lista de modelos
// db.FileModel.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo a configuração do banco de dados e o
// zap.Logger.
//
// - categId: identificador único da categoria cujos arquivos devem ser
// recuperados.
//
// Retorno:
//
// - []db.FileModel: uma lista de arquivos contendo os campos FileId, CategId,
// Name, Extension, Mimetype e UpdatedAt.
//
// - error: um erro é retornado caso a consulta ou o processamento dos
// resultados falhe.
func QueryAllFiles(
	ctx *context.Context,
	categId uuid.UUID,
) ([]db.FileModel, error) {
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
	rows, err := ctx.DB.Query(query, categId.String())
	if err != nil {
		return files, fmt.Errorf(
			"não foi possível realizar query de todos os arquivos: %s",
			err,
		)
	}
	defer closeRows(ctx, rows)

	// Iterar por cada uma das linhas
	for rows.Next() {
		var f db.FileModel
		err := rows.Scan(
			&f.FileId,
			&f.CategId,
			&f.Name,
			&f.Extension,
			&f.Mimetype,
			&f.UpdatedAt,
		)
		if err != nil {
			return files, fmt.Errorf(
				"não foi possível obter todas os arquivos: %v",
				err,
			)
		}
		files = append(files, f)
	}

	return files, nil
}

// QueryUserById realiza uma consulta ao banco de dados para buscar um usuário
// pelo ID.
//
// A função executa uma query SQL para recuperar os campos UserId, Name e
// UpdatedAt de um usuário específico, identificado pelo userId fornecido.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração e banco de
// dados.
//
// - userId: o uuid.UUID do usuário a ser consultado.
//
// Retorno:
//
// - db.UserModel: estrutura contendo os dados do usuário encontrado.
//
// - error: se ocorrer algum erro durante a consulta, um erro é retornado.
func QueryUserById(
	ctx *context.Context,
	userId uuid.UUID,
) (db.UserModel, error) {
	var user db.UserModel

	// Query
	schema := &ctx.Config.Database.Schema
	query := fmt.Sprintf(
		`SELECT %s,%s,%s,%s
		FROM %s.%s
		WHERE %s = :user_id`,
		schema.UserTable.Columns.UserId,
		schema.UserTable.Columns.Name,
		schema.UserTable.Columns.Password,
		schema.UserTable.Columns.UpdatedAt,
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.UserId,
	)

	// Obtenção da linha
	err := ctx.DB.QueryRow(query, userId.String()).Scan(
		&user.UserId,
		&user.Name,
		&user.Password,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, fmt.Errorf("não foi possível obter usuário %s: %s", userId, err)
	}

	return user, nil
}

// QueryCategoryById realiza uma consulta ao banco de dados para buscar uma
// categoria pelo seu ID.
//
// A função executa uma query SQL para recuperar os campos CategId, UserId,
// Name e UpdatedAt de uma categoria específica, identificada pelo categId
// fornecido.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração e banco de
// dados.
//
// - categId: o uuid.UUID da categoria a ser consultada.
//
// Retorno:
//
// - db.CategModel: estrutura contendo os dados da categoria encontrada.
//
// - error: se ocorrer algum erro durante a consulta, um erro é retornado.
func QueryCategoryById(
	ctx *context.Context,
	categId uuid.UUID,
) (db.CategModel, error) {
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
	err := ctx.DB.QueryRow(query, categId.String()).Scan(
		&categ.CategId,
		&categ.UserId,
		&categ.Name,
		&categ.UpdatedAt,
	)
	if err != nil {
		return categ, fmt.Errorf("não foi possível obter categoria %s: %s", categId, err)
	}

	return categ, nil
}

// QueryFileById realiza uma consulta ao banco de dados para buscar um arquivo
// pelo seu ID.
//
// A função executa uma query SQL para recuperar os campos FileId, CategId,
// Name, Extension, Mimetype e UpdatedAt de um arquivo específico, identificado
// pelo fileId fornecido.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração e banco de
// dados.
//
// - fileId: o uuid.UUID do arquivo a ser consultado.
//
// Retorno:
//
// - db.FileModel: estrutura contendo os dados do arquivo encontrado.
//
// - error: se ocorrer algum erro durante a consulta, um erro é retornado.
func QueryFileById(
	ctx *context.Context,
	fileId uuid.UUID,
) (db.FileModel, error) {
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
	err := ctx.DB.QueryRow(query, fileId.String()).Scan(
		&file.FileId,
		&file.CategId,
		&file.Name,
		&file.Extension,
		&file.Mimetype,
		&file.UpdatedAt,
	)
	if err != nil {
		return file, fmt.Errorf("não foi possível obter arquivo %s: %s", fileId, err)
	}

	return file, nil
}

package store

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/types/db"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// rollbackCreate realiza um rollback em caso de falha durante a criação de
// uma entidade no banco de dados ou no sistema de arquivos.
//
// Parâmetros:
//
// - ctx: contexto da aplicação contendo zap.Logger, configuração e recursos
// compartilhados.
//
// - rbData: estrutura CreateRollbackData.
func rollbackCreate(ctx *context.Context, rbData CreateRollbackData) {
	if rbData.Tx != nil && (rbData.DB != nil || rbData.FS != nil) {
		// Tentativas de rollback no banco
		if rbData.DB != nil {
			err := rbData.Tx.Rollback()
			if err != nil {
				ctx.Logger.Error(
					"Tentativa de rollback falhou",
					zap.Error(err),
				)
			}
		}

		// Tentativa de exclusão do arquivo/diretório, caso tenha sido criado
		if rbData.Path != "" && rbData.FS != nil {
			err := ctx.FileSystem.DeleteEntity(rbData.Path)
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
// - p: estrutura UserParams contendo as informações do usuário a ser criado.
//
// Retorno:
//
// - error: retorna um erro caso ocorra qualquer falha durante o processo de
// criação, incluindo erros de transação no banco de dados ou no sistema de
// arquivos.
func CreateUser(ctx *context.Context, p *UserParams) error {
	rbErrors := &RollbackErrors{}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	userId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível criar UUID: %v", err)
	}
	path := filepath.Join(ctx.FileSystem.Root, userId.String())

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível transação: %v", rbErrors.DB)
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

	// Criação
	_, rbErrors.DB = tx.Exec(
		insert,
		sql.Named("user_id", userId.String()),
		sql.Named("name", p.Name),
		sql.Named("password", p.Password),
		sql.Named("updated_at", ts),
	)
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível criar usuário: %v", rbErrors.DB)
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.CreateEntity(path, nil, fs.User)
	if rbErrors.FS != nil {
		return rbErrors.FS
	}

	// Confirmar a transação no banco
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", rbErrors.DB)
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
// - p: estrutura CategParams contendo as informações da categoria a ser
// criada.
//
// Retorno:
//
// - error: retorna um erro caso ocorra qualquer falha durante o processo de
// criação, incluindo erros de transação no banco de dados ou no sistema de
// arquivos.
func CreateCategory(ctx *context.Context, p *CategParams) error {
	rbErrors := &RollbackErrors{}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	categId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível criar UUID: %v", err)
	}
	path := filepath.Join(ctx.FileSystem.Root, p.UserId.String(), categId.String())

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível transação: %v", rbErrors.DB)
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
		return fmt.Errorf("não foi possível criar categoria: %v", rbErrors.DB)
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.CreateEntity(path, nil, fs.Category)
	if rbErrors.FS != nil {
		return rbErrors.FS
	}

	// Confirmar a transação no banco
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", rbErrors.DB)
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
// - p: estrutura FileParams contendo as informações do arquivo a ser criado.
//
// Retorno:
//
// - error: retorna um erro caso ocorra falha durante o processo de criação no
// banco de dados ou no sistema de arquivos.
func CreateFile(ctx *context.Context, p *FileParams) error {
	rbErrors := &RollbackErrors{}

	// Geração do UUID e Timestamp
	ts := time.Now().Unix()
	fileId, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("não foi possível criar UUID: %v", err)
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
		return fmt.Errorf("não foi possível transação: %v", rbErrors.DB)
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
		return fmt.Errorf("não foi possível criar arquivo: %v", rbErrors.DB)
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.CreateEntity(path, p.Content, fs.File)
	if rbErrors.FS != nil {
		return rbErrors.FS
	}

	// Confirmar a transação
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível confirmar transação: %v", rbErrors.DB)
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
	rows, err := ctx.DB.Query(query, sql.Named("user_id", userId.String()))
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
	rows, err := ctx.DB.Query(query, sql.Named("categ_id", categId.String()))
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
	err := ctx.DB.QueryRow(query, sql.Named("user_id", userId.String())).Scan(
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
	err := ctx.DB.QueryRow(query, sql.Named("categ_id", categId.String())).Scan(
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
	err := ctx.DB.QueryRow(query, sql.Named("file_id", fileId.String())).Scan(
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

// rollbackUpdate realiza um rollback em caso de falha durante a atualização
// de uma entidade no banco de dados ou no sistema de arquivos.
//
// Parâmetros:
//
//   - ctx: contexto da aplicação contendo zap.Logger, configuração e recursos
//     compartilhados.
//
//   - rbData: estrutura UpdateRollbackData.
func rollbackUpdate(ctx *context.Context, rbData UpdateRollbackData) {
	if rbData.Tx != nil && (rbData.DB != nil || rbData.FS != nil) {
		// Tentativas de rollback no banco
		if rbData.DB != nil {
			err := rbData.Tx.Rollback()
			if err != nil {
				ctx.Logger.Error(
					"Tentativa de rollback falhou",
					zap.Error(err),
				)
			}
		}

		// Tentativa de mover arquivo para o caminho original, caso tenha sido
		// movido
		if rbData.OldPath != "" && rbData.NewPath != "" && rbData.FS != nil {
			err := ctx.FileSystem.UpdateEntity(rbData.NewPath, rbData.OldPath)
			if err != nil {
				ctx.Logger.Error(
					"Tentativa de rollback falhou",
					zap.Error(err),
				)
			}
		}
	}
}

// UpdateUser atualiza os dados de um usuário no banco de dados.
//
// A função utiliza uma transação para garantir que as alterações sejam
// aplicadas de forma atômica. Apenas os campos fornecidos em p (como Name
// ou Password) serão atualizados, juntamente com o campo UpdatedAt que é
// automaticamente gerado com o timestamp atual.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração, logger
// e conexão com o banco de dados.
//
// - p: estrutura do tipo UserUpdate contendo os dados a serem atualizados
// para o usuário. Campos vazios em p serão ignorados.
//
// Retorno:
//
// - error: se ocorrer algum erro durante o processo de atualização, incluindo
// erros ao iniciar ou confirmar a transação, ou ao executar a query de
// atualização.
func UpdateUser(ctx *context.Context, p UserUpdate) error {
	rbErrors := &RollbackErrors{}

	// Geração do Timestamp
	ts := time.Now().Unix()

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", rbErrors.DB)
	}

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
	var args []sql.NamedArg
	var set []string
	if p.Name != p.OldName {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.UserTable.Columns.Name+" = :name")
	}
	if p.Password != "" {
		args = append(args, sql.Named("password", p.Password))
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
	res, rbErrors.DB = tx.Exec(update, args)
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível atualizar usuário: %s", rbErrors.DB)
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		return rbErrors.DB
	}

	// Confirmar a transação
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", rbErrors.DB)
	}

	return nil
}

// UpdateCategory atualiza os dados de uma categoria no banco de dados e realiza
// alterações correspondentes no sistema de arquivos, se necessário.
//
// A função utiliza uma transação para garantir a atomicidade das operações,
// tanto no banco de dados quanto no sistema de arquivos. Os parâmetros que
// serão atualizados (como UserId ou Name) são verificados dinamicamente.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração,
// zap.logger, conexão com o banco de dados e o sistema de arquivos.
//
// - p: estrutura do tipo CategUpdate contendo os dados da categoria a serem
// atualizados. Inclui o ID do usuário atual, o novo ID do usuário (se for
// alterado), e o nome da categoria atual e novo.
//
// Retorno:
//
// - error: retorna um erro se ocorrerem falhas em qualquer etapa do processo,
// como falha ao iniciar ou confirmar a transação, ou erro ao atualizar o
// sistema de arquivos.
func UpdateCategory(ctx *context.Context, p CategUpdate) error {
	rbErrors := &RollbackErrors{}

	// Geração do Timestamp e caminhos
	ts := time.Now().Unix()
	oldPath := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
	)
	newPath := ctx.FileSystem.Root

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []sql.NamedArg
	var set []string
	if p.UserId != p.OldUserId {
		newPath = filepath.Join(newPath, p.UserId.String())
		if !ctx.FileSystem.EntityExists(newPath) {
			rbErrors.DB = fmt.Errorf(
				"não pôde atualizar categoria: diretório do usuário %s não existe",
				p.OldUserId.String(),
			)
			return rbErrors.DB
		}
		args = append(args, sql.Named("user_id", p.UserId))
		set = append(set, schema.CategTable.Columns.UserId+" = :user_id")
	} else {
		newPath = filepath.Join(newPath, p.OldUserId.String())
	}
	newPath = filepath.Join(newPath, p.CategId.String())
	if p.Name != p.OldName {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.CategTable.Columns.Name+" = :name")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.CategTable.Columns.UpdatedAt+" = :updated_at")

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", rbErrors.DB)
	}

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
	res, rbErrors.DB = tx.Exec(update, args)
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível atualizar categoria: %s", rbErrors.DB)
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
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
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", rbErrors.DB)
	}

	return nil
}

// UpdateFile atualiza os dados de um arquivo no banco de dados e realiza
// alterações correspondentes no sistema de arquivos, se necessário.
//
// A função utiliza uma transação para garantir a atomicidade das operações,
// tanto no banco de dados quanto no sistema de arquivos. Os parâmetros que
// serão atualizados (como CategId, Name, Extension ou Mimetype) são verificados
// dinamicamente.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração,
// zap.logger, conexão com o banco de dados e o sistema de arquivos.
//
// - p: estrutura do tipo FileUpdate contendo os dados do arquivo a serem
// atualizados. Inclui IDs do usuário e da categoria, nome atual e novo,
// extensão, mimetype e conteúdo, se aplicável.
//
// Retorno:
//
// - error: retorna um erro se ocorrerem falhas em qualquer etapa do processo,
// como falha ao iniciar ou confirmar a transação, ou erro ao atualizar o
// sistema de arquivos.
func UpdateFile(ctx *context.Context, p FileUpdate) error {
	rbErrors := &RollbackErrors{}

	// Geração do Timestamp e caminhos
	ts := time.Now().Unix()
	oldPath := filepath.Join(
		ctx.FileSystem.Root,
		p.UserId.String(),
		p.CategId.String(),
		p.FileId.String()+p.Extension,
	)
	newPath := filepath.Join(ctx.FileSystem.Root, p.UserId.String())

	// Checagem dos parâmetros a serem atualizados
	schema := &ctx.Config.Database.Schema
	var args []sql.NamedArg
	var set []string
	if p.CategId != p.OldCategId {
		newPath = filepath.Join(newPath, p.CategId.String())
		if !ctx.FileSystem.EntityExists(newPath) {
			rbErrors.DB = fmt.Errorf(
				"não pôde atualizar arquivo: diretório da categoria %s "+
					"não existe para o usuário %s",
				p.CategId.String(),
				p.UserId.String(),
			)
			return rbErrors.DB
		}
		args = append(args, sql.Named("categ_id", p.CategId))
		set = append(set, schema.FileTable.Columns.CategId+" = :categ_id")
	} else {
		newPath = filepath.Join(newPath, p.OldCategId.String())
	}
	newPath = filepath.Join(newPath, p.FileId.String())
	if p.Name != p.OldName {
		args = append(args, sql.Named("name", p.Name))
		set = append(set, schema.FileTable.Columns.Name+" = :name")
	}
	if p.Extension != p.OldExtension {
		args = append(args, sql.Named("extension", p.Extension))
		set = append(set, schema.FileTable.Columns.Extension+" = :extension")
		newPath = newPath + p.Extension
	}
	if p.Mimetype != p.OldMimetype {
		args = append(args, sql.Named("mimetype", p.Mimetype))
		set = append(set, schema.FileTable.Columns.Mimetype+" = :mimetype")
	}
	args = append(args, sql.Named("updated_at", ts))
	set = append(set, schema.CategTable.Columns.UpdatedAt+" = :updated_at")

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", rbErrors.DB)
	}

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
	res, rbErrors.DB = tx.Exec(update, args)
	if rbErrors.DB != nil {
		return fmt.Errorf("não foi possível atualizar arquivo: %s", rbErrors.DB)
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
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
		defer func() {
			if rbErrors.FS != nil {
				// Tentativa de exclusão do arquivo criado, caso tenha erros
				err := ctx.FileSystem.DeleteEntity(newPath)
				if err != nil {
					ctx.Logger.Error(
						"Tentativa de limpeza falhou",
						zap.Error(err),
					)
				}
			} else {
				// Tentativa de exclusão do arquivo antigo, quando não tiver
				// erros
				err := ctx.FileSystem.DeleteEntity(oldPath)
				if err != nil {
					ctx.Logger.Error(
						"Tentativa de limpeza falhou",
						zap.Error(err),
					)
				}
			}
		}()
	}

	// Confirmar a transação
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", rbErrors.DB)
	}

	return nil
}

// rollbackDelete realiza o rollback de uma operação de exclusão no banco de
// dados e/ou no sistema de arquivos.
//
// A função é chamada para desfazer alterações realizadas durante uma operação
// de exclusão, em caso de falha. Ela tenta reverter mudanças no banco de dados
// e recriar arquivos no sistema de arquivos com base nos dados fornecidos.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração, banco de
// dados e sistema de arquivos.
//
// - rbData: estrutura DeleteRollbackData.
func rollbackDelete(ctx *context.Context, rbData DeleteRollbackData) {
	if rbData.Tx != nil && (rbData.DB != nil || rbData.FS != nil) {
		// Tentativas de rollback no banco
		if rbData.DB != nil {
			err := rbData.Tx.Rollback()
			if err != nil {
				ctx.Logger.Error(
					"Tentativa de rollback falhou",
					zap.Error(err),
				)
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
				ctx.Logger.Error(
					"Tentativa de rollback falhou",
					zap.Error(err),
				)
			}
		}
	}
}

// DeleteUser realiza a exclusão de um usuário no banco de dados e no sistema de
// arquivos.
//
// A função é responsável por remover o registro do usuário identificado pelo
// UserId no banco de dados, bem como excluir o diretório associado a esse
// usuário no sistema e arquivos. Em caso de falha, um rollback é realizado para
// desfazer alterações.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração, banco de
// dados e sistema de arquivos.
//
// - p: estrutura UserDelete.
//
// Retorno:
//
// - error: retorna um erro se ocorrer alguma falha durante o processo de
// exclusão.
func DeleteUser(ctx *context.Context, p UserDelete) error {
	rbErrors := &RollbackErrors{}

	// Caminho
	path := filepath.Join(ctx.FileSystem.Root, p.UserId.String())

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", rbErrors.DB)
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
		return fmt.Errorf("não foi possível excluir usuário: %s", rbErrors.DB)
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		return rbErrors.DB
	}

	// Confirmar a transação
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", rbErrors.DB)
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.DeleteEntity(path)
	if rbErrors.FS != nil {
		return rbErrors.FS
	}

	return nil
}

// DeleteCategory realiza a exclusão de uma categoria no banco de dados e no
// sistema de arquivos.
//
// A função é responsável por remover o registro da categoria identificado pelo
// CategId no banco de dados, bem como excluir o diretório associado a essa
// categoria no sistema de arquivos. Em caso de falha, um rollback é realizado
// para desfazer alterações.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração, banco de
// dados e sistema de arquivos.
//
// - p: estrutura CategDelete.
//
// Retorno:
//
// - error: retorna um erro se ocorrer alguma falha durante o processo de
// exclusão.
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
		return fmt.Errorf("erro ao iniciar transação: %v", rbErrors.DB)
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
		return fmt.Errorf("não foi possível excluir categoria: %s", rbErrors.DB)
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		return rbErrors.DB
	}

	// Confirmar a transação
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", rbErrors.DB)
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.DeleteEntity(path)
	if rbErrors.FS != nil {
		return rbErrors.FS
	}

	return nil
}

// DeleteFile realiza a exclusão de um arquivo no banco de dados e no sistema de
// arquivos.
//
// A função é responsável por remover o registro do arquivo identificado pelo
// FileId no banco de dados, bem como excluir o arquivo associado no sistema de
// arquivos. Antes de excluir o arquivo, seu conteúdo é lido e armazenado como
// backup para uso em caso de rollback. Em caso de falha, um rollback é
// realizado para desfazer alterações.
//
// Parâmetros:
//
// - ctx: contexto da aplicação, contendo informações de configuração, banco de
// dados e sistema de arquivos.
//
// - p: estrutura FileDelete, contendo os dados necessários para exclusão.
//
// Retorno:
//
// - error: retorna um erro se ocorrer alguma falha durante o processo de
// exclusão.
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
		rbErrors.FS = fmt.Errorf("erro ao abrir %s: %v", path, err)
		return rbErrors.FS
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ctx.Logger.Error(
				"Erro ao fechar arquivo para leitura",
				zap.Error(err),
			)
		}
	}(file)
	// Leitura
	var backupContent []byte
	backupContent, rbErrors.FS = io.ReadAll(file)
	if rbErrors.FS != nil {
		rbErrors.FS = fmt.Errorf("erro ao ler %s para backup: %v", path, rbErrors.FS)
		return rbErrors.FS
	}

	// Iniciar uma transação
	var tx *sql.Tx
	tx, rbErrors.DB = ctx.DB.Begin()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", rbErrors.DB)
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
		return fmt.Errorf("não foi possível excluir arquivo: %s", rbErrors.DB)
	} else if n, _ := res.RowsAffected(); n > 1 {
		rbErrors.DB = fmt.Errorf("mais de uma linha afetada")
		return rbErrors.DB
	}

	// Confirmar a transação
	rbErrors.DB = tx.Commit()
	if rbErrors.DB != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", rbErrors.DB)
	}

	// Transação no sistema de arquivos
	rbErrors.FS = ctx.FileSystem.DeleteEntity(path)
	if rbErrors.FS != nil {
		return rbErrors.FS
	}

	return nil
}

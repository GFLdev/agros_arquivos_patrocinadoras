package db

import (
	"agros_arquivos_patrocinadoras/filerepo/services/logger"
	db2 "agros_arquivos_patrocinadoras/pkg/types/db"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func QueryAllUsers(db *sql.DB) ([]db2.UserModel, error) {
	var users []db2.UserModel

	query := `SELECT uuid_patroc
		, nome_patroc
		, ts_modificado
		FROM patrocinadoras_repositorio`

	rows, err := db.Query(query)
	if err != nil {
		return users, fmt.Errorf(
			"nao foi possivel realizar query de todos os usuários: %s",
			err,
		)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.CreateLogger().Warn(
				"Erro ao fechar linhas da query de usuários",
				zap.Error(err),
			)
		}
	}(rows)

	// Iterar por cada linha
	for rows.Next() {
		var u db2.UserModel
		err := rows.Scan(&u.UserId, &u.Name, &u.UpdatedAt)
		if err != nil {
			return users, fmt.Errorf(
				"não foi possível obter todos os usuários: %s",
				err,
			)
		}
		u.Password = "OMITIDO"
		users = append(users, u)
	}

	return users, nil
}

func QueryUserById(db *sql.DB, userId uuid.UUID) (db2.UserModel, error) {
	var user db2.UserModel

	query := `SELECT uuid_patroc
		, nome_patroc
		, ts_modificado
		FROM patrocinadoras_repositorio
		WHERE uuid_patroc = $1`

	err := db.QueryRow(query, userId.String()).Scan(
		&user.UserId,
		&user.Name,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, fmt.Errorf("não foi possível obter usuário: %s", err)
	}

	return user, nil
}

func QueryAllCategories(db *sql.DB, userId uuid.UUID) ([]db2.CategoryModel, error) {
	var categs []db2.CategoryModel

	query := `SELECT uuid_categ
     	, uuid_patroc
		, nome_categ
		, ts_modificado
		FROM patrocinadoras_categoria
		WHERE uuid_patroc = $1`

	rows, err := db.Query(query, userId.String())
	if err != nil {
		return categs, fmt.Errorf(
			"nao foi possivel realizar query de todas as categorias: %s",
			err,
		)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.CreateLogger().Warn(
				"Erro ao fechar linhas da query de categorias",
				zap.Error(err),
			)
		}
	}(rows)

	// Iterar por cada linha
	for rows.Next() {
		var c db2.CategoryModel
		err := rows.Scan(&c.CategId, &c.UserId, &c.Name, &c.UpdatedAt)
		if err != nil {
			return categs, fmt.Errorf(
				"não foi possível obter todas as categorias: %s",
				err,
			)
		}
		categs = append(categs, c)
	}

	return categs, nil
}

func QueryCategoryById(db *sql.DB, categId uuid.UUID) (db2.CategoryModel, error) {
	var categ db2.CategoryModel

	query := `SELECT uuid_categ
     	, uuid_patroc
		, nome_categ
		, ts_modificado
		FROM patrocinadoras_categoria
		WHERE uuid_categ = $1`

	err := db.QueryRow(query, categId.String()).Scan(
		&categ.CategId,
		&categ.UserId,
		&categ.Name,
		&categ.UpdatedAt,
	)
	if err != nil {
		return categ, fmt.Errorf("não foi possível obter categoria: %s", err)
	}

	return categ, nil
}

func QueryAllFiles(db *sql.DB, categId uuid.UUID) ([]db2.FileModel, error) {
	var files []db2.FileModel

	query := `SELECT f.uuid_arquivo
		, f.uuid_categ
		, f.nome_arquivo
		, f.extensao
		, f.mimetype
		, f.ts_modificado
		FROM patrocinadoras_arquivo f
		WHERE f.uuid_categ = $1`

	rows, err := db.Query(query, categId.String())
	if err != nil {
		return files, fmt.Errorf(
			"nao foi possivel realizar query de todos os arquivos: %s",
			err,
		)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.CreateLogger().Warn(
				"Erro ao fechar linhas da query de arquivos",
				zap.Error(err),
			)
		}
	}(rows)

	// Iterar por cada linha
	for rows.Next() {
		var f db2.FileModel
		err := rows.Scan(&f.FileId, &f.CategId, &f.Name, &f.Extension, &f.Mimetype, &f.UpdatedAt)
		if err != nil {
			return files, fmt.Errorf(
				"não foi possível obter todos os arquivos: %s",
				err,
			)
		}
		files = append(files, f)
	}

	return files, nil
}

func QueryFileById(db *sql.DB, fileId uuid.UUID) (db2.FileModel, error) {
	var file db2.FileModel

	query := `SELECT a.uuid_arquivo
		, a.uuid_categ
		, a.nome_arquivo
		, a.extensao
		, a.mimetype
		, a.ts_modificado
		FROM patrocinadoras_arquivo a
		WHERE a.uuid_arquivo = $1`

	err := db.QueryRow(query, fileId.String()).Scan(
		&file.FileId,
		&file.CategId,
		&file.Name,
		&file.Extension,
		&file.Mimetype,
		&file.UpdatedAt,
	)
	if err != nil {
		return file, fmt.Errorf("não foi possível obter arquivo: %s", err)
	}

	return file, nil
}

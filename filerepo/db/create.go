package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func CreateUser(
	db *sql.DB,
	userId uuid.UUID,
	name string,
	password string,
) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Criação
	insert := `INSERT INTO patrocinadoras_repositorio
  		(uuid_patroc, nome_patroc, senha, ts_modificado)
		VALUES ($1, $2, $3, $4)`
	ts := time.Now().Unix()

	_, err = tx.Exec(insert, userId.String(), name, password, ts)
	if err != nil {
		return fmt.Errorf("não foi possível criar usuário: %v", err)
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

func CreateCategory(
	db *sql.DB,
	categId uuid.UUID,
	userId uuid.UUID,
	name string,
) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}

	// Criação
	insert := `INSERT INTO patrocinadoras_categoria
  		(uuid_categ, uuid_patroc, nome_categ, ts_modificado)
		VALUES ($1, $2, $3, $4)`
	ts := time.Now().Unix()

	_, err = tx.Exec(insert, categId.String(), userId.String(), name, ts)
	if err != nil {
		return fmt.Errorf("não foi possível criar categoria: %v", err)
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

func CreateFile(
	db *sql.DB,
	fileId uuid.UUID,
	categId uuid.UUID,
	name string,
	extension string,
	mimeType string,
) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Criação
	insert := `INSERT INTO patrocinadoras_arquivo
  		(uuid_arquivo, uuid_categ, nome_arquivo, extensao, mimetype, ts_modificado)
		VALUES ($1, $2, $3, $4, $5, $6)`
	ts := time.Now().Unix()

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

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

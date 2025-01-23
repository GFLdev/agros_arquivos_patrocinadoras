package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

func DeleteUser(db *sql.DB, userId uuid.UUID) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Exclusão
	del := `DELETE FROM patrocinadoras_repositorio
       WHERE uuid_patroc = $1`

	res, err := tx.Exec(del, userId)
	if err != nil {
		return fmt.Errorf("não foi possível excluir usuario: %s", err)
	} else if n, _ := res.RowsAffected(); n > 1 {
		return fmt.Errorf("mais de uma linha afetada")
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

func DeleteCategory(db *sql.DB, categId uuid.UUID) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Exclusão
	del := `DELETE FROM patrocinadoras_categoria
       WHERE uuid_categ = $1`

	res, err := tx.Exec(del, categId)
	if err != nil {
		return fmt.Errorf("não foi possível excluir categoria: %s", err)
	} else if n, _ := res.RowsAffected(); n > 1 {
		return fmt.Errorf("mais de uma linha afetada")
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

func DeleteFile(db *sql.DB, fileId uuid.UUID) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Exclusão
	del := `DELETE FROM patrocinadoras_arquivo
		WHERE uuid_arquivo = $1`

	res, err := tx.Exec(del, fileId)
	if err != nil {
		return fmt.Errorf("não foi possível excluir arquivo: %s", err)
	} else if n, _ := res.RowsAffected(); n > 1 {
		return fmt.Errorf("mais de uma linha afetada")
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

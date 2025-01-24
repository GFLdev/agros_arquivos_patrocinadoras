package db

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func UpdateUser(
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

	// Checagem dos parâmetros a serem atualizados
	var args []interface{}
	var set []string
	if name != "" {
		args = append(args, name)
		set = append(set, "nome_patroc = ?")
	}
	if password != "" {
		args = append(args, password)
		set = append(set, "senha = ?")
	}

	ts := time.Now().Unix()
	args = append(args, ts, userId.String())
	set = append(set, "ts_modificado = ?")

	// Atualização
	update := fmt.Sprintf(`UPDATE patrocinadoras_repositorio
				SET %s
				WHERE uuid_patroc = ?`, strings.Join(set, ","))

	res, err := tx.Exec(update, args...)
	if err != nil {
		return fmt.Errorf("não foi possível atualizar categoria: %s", err)
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

func UpdateCategory(
	db *sql.DB,
	categId uuid.UUID,
	userId string,
	name string,
) error {
	// Iniciar uma transação
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Checagem dos parâmetros a serem atualizados
	var args []any
	var set []string
	if _, err = uuid.Parse(userId); err == nil {
		args = append(args, userId)
		set = append(set, "uuid_patroc = ?")
	}
	if name != "" {
		args = append(args, name)
		set = append(set, "nome_categ = ?")
	}

	ts := time.Now().Unix()
	args = append(args, ts, categId.String())
	set = append(set, "ts_modificado = ?")

	// Atualização
	update := fmt.Sprintf(`UPDATE patrocinadoras_categoria
				SET %s
				WHERE uuid_categ = ?`,
		strings.Join(set, ","),
	)

	res, err := tx.Exec(update, args...)
	if err != nil {
		return fmt.Errorf("não foi possível atualizar categoria: %s", err)
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

func UpdateFile(
	db *sql.DB,
	fileId uuid.UUID,
	categId string,
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

	// Checagem dos parâmetros a serem atualizados
	var args []any
	var set []string
	if _, err = uuid.Parse(categId); err == nil {
		args = append(args, categId)
		set = append(set, "uuid_categ = ?")
	}
	if name != "" {
		args = append(args, name)
		set = append(set, "nome_arquivo = ?")
	}
	if extension != "" {
		args = append(args, extension)
		set = append(set, "extensao = ?")
	}
	if mimeType != "" {
		args = append(args, mimeType)
		set = append(set, "mimetype = ?")
	}

	ts := time.Now().Unix()
	args = append(args, ts, fileId.String())
	set = append(set, "ts_modificado = ?")

	// Atualização
	update := fmt.Sprintf(`UPDATE patrocinadoras_arquivo
				SET %s
				WHERE uuid_arquivo = ?`,
		strings.Join(set, ","),
	)

	res, err := tx.Exec(update, args...)
	if err != nil {
		return fmt.Errorf("não foi possível atualizar arquivo: %s", err)
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

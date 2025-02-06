package test

import (
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"database/sql"
	"fmt"
	"os"
	"testing"
)

func newContext() *context.Context {
	logr := logger.CreateLogger()
	cfg := config.LoadConfig(logr)
	dataBase := db.GetSqlDB(&cfg.Database, logr)
	return &context.Context{
		Logger: logr,
		Config: cfg,
		DB:     dataBase,
	}
}

func TestMain(m *testing.M) {
	// Executa os testes
	code := m.Run()

	// Cleanup e finalização
	ctx := newContext()
	defer func(DB *sql.DB) {
		_ = DB.Close()
	}(ctx.DB)

	// Apagar dados do banco
	schema := ctx.Config.Database.Schema
	delUsers := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE %s <> :adminName",
		schema.Name,
		schema.UserTable.Name,
		schema.UserTable.Columns.Name,
	)
	delCategories := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE 0 = 0",
		schema.Name,
		schema.CategTable.Name,
	)
	delFiles := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE 0 = 0",
		schema.Name,
		schema.FileTable.Name,
	)

	tx, err := ctx.DB.Begin()
	if err != nil {
		panic(err)
	}
	_, _ = tx.Exec(delFiles)
	_, _ = tx.Exec(delCategories)
	_, _ = tx.Exec(delUsers, sql.Named("adminName", ctx.Config.AdminName))
	_ = tx.Commit()

	// Finaliza os testes
	os.Exit(code)
}

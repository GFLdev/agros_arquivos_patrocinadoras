package test

import (
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"fmt"
	"os"
	"testing"
)

const (
	TestRoot = "test_data"
)

func newContext() *context.Context {
	logr := logger.CreateLogger()
	cfg := config.LoadConfig(logr)
	filesystem := &fs.FileSystem{Root: TestRoot}
	dataBase := db.GetSqlDB(&cfg.Database, logr)
	return &context.Context{
		Logger:     logr,
		Config:     cfg,
		FileSystem: filesystem,
		DB:         dataBase,
	}
}

func TestMain(m *testing.M) {
	// Configuração inicial antes de todos os testes
	if err := os.MkdirAll(TestRoot, os.ModePerm); err != nil {
		panic(err)
	}

	// Executa os testes
	code := m.Run()

	// Cleanup e finalização
	_ = os.RemoveAll(TestRoot)
	ctx := newContext()

	// Apagar dados do banco
	schema := ctx.Config.Database.Schema
	delUsers := fmt.Sprintf(
		"DELETE FROM %s.%s WHERE 0 = 0",
		schema.Name,
		schema.UserTable.Name,
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
	_, _ = tx.Exec(delUsers)
	_ = tx.Commit()

	// Finaliza os testes
	os.Exit(code)
}

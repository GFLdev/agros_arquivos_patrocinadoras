package test

import (
	"os"
	"testing"
)

const (
	TestRoot = "test_data"
)

func TestMain(m *testing.M) {
	// Configuração inicial antes de todos os testes
	if err := os.MkdirAll(TestRoot, os.ModePerm); err != nil {
		panic(err)
	}

	// Executa os testes
	code := m.Run()

	// Cleanup e finalização
	_ = os.RemoveAll(TestRoot)

	// Finaliza os testes
	os.Exit(code)
}

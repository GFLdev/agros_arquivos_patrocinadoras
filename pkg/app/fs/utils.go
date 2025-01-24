package fs

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
)

func WriteToFile(path string, data []byte, logr *zap.Logger) error {
	// Abertura do arquivo
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro a abrir %s: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logr.Fatal("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	// Criação de um buffer de escrita
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			logr.Fatal("Erro ao liberar buffer de escrita", zap.Error(err))
		}
	}(writer)

	// Escrever todo o slice
	for len(data) > 0 {
		n, err := writer.Write(data)
		if err != nil {
			return fmt.Errorf("erro ao escrever em %s: %v", path, err)
		}
		data = data[n:]
	}

	return nil
}

func DeleteFile(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("erro ao excluir %s: %v", path, err)
	}
	return nil
}

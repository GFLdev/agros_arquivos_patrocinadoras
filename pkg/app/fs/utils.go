package fs

import (
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
)

// WriteToFile cria ou abre um arquivo no caminho especificado e escreve
// os dados nele.
//
// Parâmetros:
//   - path: caminho do arquivo a ser criado ou sobrescrito.
//   - data: dados em formato de slice de bytes que serão escritos no arquivo.
//
// Retorno:
//   - error: retorna um erro caso ocorra algum problema, ou nil caso contrário.
func WriteToFile(path string, data []byte) error {
	logr := logger.CreateLogger()

	// Abertura do arquivo
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro a abrir %s: %w", path, err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logr.Error("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	// Criação de um buffer de escrita
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		if err := writer.Flush(); err != nil {
			logr.Error("Erro ao liberar buffer de escrita", zap.Error(err))
		}
	}(writer)

	// Escrever todo o slice
	for len(data) > 0 {
		n, err := writer.Write(data)
		if err != nil {
			return fmt.Errorf("erro ao escrever em %s: %w", path, err)
		}
		data = data[n:]
	}

	return nil
}

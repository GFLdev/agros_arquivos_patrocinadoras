package fs

import (
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
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

	// Criação de um buffer de escrita
	writer := bufio.NewWriter(file)
	defer func(file *os.File, writer *bufio.Writer) {
		if err = writer.Flush(); err != nil {
			logr.Error("Erro ao liberar buffer de escrita", zap.Error(err))
		}
		if err = file.Close(); err != nil {
			logr.Error("Erro ao fechar "+path, zap.Error(err))
		}
	}(file, writer)

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

// EntityIsEmpty verifica se um diretório no caminho especificado está vazio.
//
// Parâmetros:
//   - path: caminho do diretório a ser verificado.
//
// Retorno:
//   - bool: retorna true se o diretório estiver vazio ou false caso contrário.
//   - error: retorna um erro caso ocorra algum problema durante a verificação.
func EntityIsEmpty(path string) (bool, error) {
	logr := logger.CreateLogger()

	// Abrindo arquivo
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			logr.Error("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	// Verificação
	info, err := file.Stat()
	if !info.IsDir() {
		return true, nil
	}
	if _, err = file.Readdirnames(1); err == io.EOF {
		return true, nil
	}
	return false, nil
}

package fs

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
)

// WriteToFile cria ou abre um arquivo no caminho especificado e escreve
// os dados nele.
//
// A função cria um novo arquivo (ou substitui o existente) no caminho
// fornecido, escreve os dados byte a byte e garante que o arquivo seja
// corretamente fechado e que o buffer de escrita, seja liberado. Em caso
// de erro ao abrir, escrever ou fechar o arquivo, o erro é retornado.
//
// Parâmetros:
//
// - path: caminho do arquivo a ser criado ou sobrescrito.
//
// - data: dados em formato de slice de bytes que serão escritos no arquivo.
//
//   - logr: zap.Logger utilizado para registrar erros ao fechar o arquivo ou
//     liberar o buffer de escrita.
//
// Retorno:
//
// - error: retorna um erro caso ocorra algum problema, ou nil caso contrário.
func WriteToFile(path string, data []byte, logr *zap.Logger) error {
	// Abertura do arquivo
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro a abrir %s: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logr.Error("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	// Criação de um buffer de escrita
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			logr.Error("Erro ao liberar buffer de escrita", zap.Error(err))
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

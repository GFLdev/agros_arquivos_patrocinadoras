package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime"
	"net/http"
	"os"
)

func GetFS(logr *zap.Logger) (*FS, error) {
	// Carregamento do arquivo de rastreamento
	repo, err := FileToStruct[FS]("repo/track.json", logr)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func GetExtension(fileType string, content *[]byte) string {
	ext, err := mime.ExtensionsByType(fileType)
	if err != nil {
		fileType = http.DetectContentType(*content)
		ext, _ = mime.ExtensionsByType(fileType)
	}

	return ext[0]
}

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

func StructToFile[T any](path string, data *T, logr *zap.Logger) error {
	// Dados JSON para escrita
	jsonPayload, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		logr.Error("Erro ao serializar JSON", zap.Error(err))
		return fmt.Errorf("erro ao serializar JSON: %v", err)
	}

	// Escrita
	err = WriteToFile(path, jsonPayload, logr)
	if err != nil {
		return err
	}

	return nil
}

func FileToStruct[T any](path string, logr *zap.Logger) (*T, error) {
	// Leitura do conteúdo em path
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir %s: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logr.Fatal("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler %s: %v", path, err)
	}

	// De-serializar JSON
	data := new(T)
	err = json.Unmarshal(byteValue, data)
	if err != nil {
		logr.Error("Erro ao de-serializar JSON", zap.Error(err))
		return nil, fmt.Errorf("erro ao de-serializar JSON: %v", err)
	}

	return data, nil
}

func DeleteFiles(path string, logr *zap.Logger) error {
	err := os.RemoveAll(path)
	if err != nil {
		logr.Error("Erro ao deletar "+path, zap.Error(err))
		return fmt.Errorf("erro ao excluir %s: %v", path, err)
	}
	return nil
}

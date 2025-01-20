package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
)

// ReadFromFile lê um arquivo e retorna seus dados como um slice de bytes.
func ReadFromFile(path string, logger *zap.Logger) ([]byte, error) {
	// Abertura do arquivo
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir %s: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	// Criação de um buffer de leitura
	reader := bufio.NewReader(file)

	// Leitura
	var contents []byte
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("erro ao ler %s: %v", path, err)
		}

		contents = append(contents, []byte(line)...)
	}

	return contents, nil
}

// WriteToFile escreve um slice de bytes em um arquivo.
func WriteToFile(path string, data []byte, logger *zap.Logger) error {
	// Abertura do arquivo
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erro a abrir %s: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal("Erro ao fechar "+path, zap.Error(err))
		}
	}(file)

	// Criação de um buffer de escrita
	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			logger.Fatal("Erro ao liberar buffer de escrita", zap.Error(err))
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

func StructToFile[T any](path string, data *T, logger *zap.Logger) error {
	// Dados JSON para escrita
	jsonPayload, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		logger.Error("Erro ao serializar JSON", zap.Error(err))
		return fmt.Errorf("erro ao serializar JSON: %v", err)
	}

	// Escrita
	err = WriteToFile(path, jsonPayload, logger)
	if err != nil {
		return err
	}

	return nil
}

func FileToStruct[T any](path string, logger *zap.Logger) (*T, error) {
	// Leitura do conteúdo em path
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir %s: %v", path, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Fatal("Erro ao fechar "+path, zap.Error(err))
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
		logger.Error("Erro ao de-serializar JSON", zap.Error(err))
		return nil, fmt.Errorf("erro ao de-serializar JSON: %v", err)
	}

	return data, nil
}

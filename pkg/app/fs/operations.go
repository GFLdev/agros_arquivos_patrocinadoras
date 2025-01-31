package fs

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"
)

// CreateEntity cria uma entidade no sistema de arquivos.
//
// Parâmetros:
//   - path: caminho onde a entidade será criada.
//   - content: conteúdo a ser escrito no arquivo (apenas para entidades do tipo
//     File). Pode ser nil para diretórios.
//   - entity: tipo da entidade a ser criada (File, User ou Category).
//
// Retorno:
//   - error: retorna um erro se ocorrer falha na criação da entidade ou se o
//     tipo de entidade for inválido.
func (fs *FileSystem) CreateEntity(path string, content *[]byte, entity EntityType) error {
	if fs.EntityExists(path) {
		return fmt.Errorf("%s já existe", path)
	}

	baseName := filepath.Base(strings.TrimSuffix(path, filepath.Ext(path)))
	if uuid.Validate(baseName) != nil {
		return fmt.Errorf("caminho inválido (UUID): %s", path)
	}

	switch entity {
	case File:
		if content == nil {
			return fmt.Errorf("conteúdo não pode ser nulo ao criar um arquivo")
		}
		return WriteToFile(path, *content)
	case User, Category:
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			var entityName string
			if entity == User {
				entityName = "usuário"
			} else {
				entityName = "categoria"
			}
			return fmt.Errorf("erro ao criar pasta da %s: %w", entityName, err)
		}
	default:
		return fmt.Errorf("tipo de entidade inválido: %d", entity)
	}
	return nil
}

// EntityExists verifica se um arquivo ou diretório existe no caminho
// especificado.
//
// Parâmetros:
//   - path: caminho do arquivo ou diretório a ser verificado.
//
// Retorno:
//   - bool: retorna true se o arquivo ou diretório existir, ou false caso
//     contrário.
func (fs *FileSystem) EntityExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// UpdateEntity atualiza a localização de uma entidade no sistema de arquivos,
// movendo-a do caminho antigo para um novo caminho especificado.
//
// Parâmetros:
//   - oldPath: string representando o caminho atual da entidade.
//   - newPath: string representando o novo caminho da entidade.
//
// Retorno:
//   - error: retorna um erro caso a operação de renomear falhe. Inclui detalhes
//     sobre o caminho antigo, o novo caminho e a causa do erro.
func (fs *FileSystem) UpdateEntity(oldPath, newPath string) error {
	// Validação
	if empty, err := EntityIsEmpty(oldPath); err != nil {
		return err
	} else if !empty {
		return fmt.Errorf("diretório %s não vazio", oldPath)
	}

	// Atualização
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("erro ao mover %s para %s: %w", oldPath, newPath, err)
	}
	return nil
}

// DeleteEntity exclui o arquivo ou diretório no caminho especificado.
//
// Parâmetros:
//   - path: caminho do arquivo ou diretório a ser excluído.
//
// Retorno:
//   - error: retorna um erro caso a exclusão falhe, ou nil caso contrário.
func (fs *FileSystem) DeleteEntity(path string) error {
	// Validação
	if empty, err := EntityIsEmpty(path); err != nil {
		return err
	} else if !empty {
		return fmt.Errorf("diretório %s não vazio", path)
	}

	// Exclusão
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("erro ao excluir %s: %w", path, err)
	}
	return nil
}

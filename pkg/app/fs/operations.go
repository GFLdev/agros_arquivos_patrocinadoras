package fs

import (
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

// CreateEntity cria uma entidade no sistema de arquivos.
// Dependendo do tipo da entidade fornecida, pode criar:
//
// - Um arquivo: escreve o conteúdo especificado no caminho fornecido.
//
// - Um diretório: cria uma pasta para usuários ou categorias.
//
// Parâmetros:
//
// - path: caminho onde a entidade será criada.
//
// - content: conteúdo a ser escrito no arquivo (apenas para entidades do tipo
// File). Pode ser nil para diretórios.
//
// - entity: tipo da entidade a ser criada (File, User ou Category).
//
// Retorno:
//
// - error: retorna um erro se ocorrer falha na criação da entidade ou se o tipo
// de entidade for inválido.
func (fs *FileSystem) CreateEntity(
	path string,
	content *[]byte,
	entity EntityType,
) error {
	logr := logger.CreateLogger()

	switch entity {
	case File:
		if content == nil {
			return fmt.Errorf("conteúdo não pode ser nulo ao criar um arquivo")
		}
		return WriteToFile(path, *content, logr)
	case User, Category:
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			var entityName string
			if entity == User {
				entityName = "usuário"
			} else if entity == Category {
				entityName = "categoria"
			}
			return fmt.Errorf("erro ao criar pasta da %s: %v", entityName, err)
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
//
// - path: caminho do arquivo ou diretório a ser verificado.
//
// Retorno:
//
// - bool: retorna true se o arquivo ou diretório existir, ou false caso
// contrário.
func (fs *FileSystem) EntityExists(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// TODO: Implementar UpdateEntity

func (fs *FileSystem) UpdateEntity() error {
	return fmt.Errorf("não implementado")
}

// DeleteEntity exclui o arquivo ou diretório no caminho especificado.
//
// Parâmetros:
//
// - path: caminho do arquivo ou diretório a ser excluído.
//
// Retorno:
//
// - error: retorna um erro caso a exclusão falhe, ou nil caso contrário.
func (fs *FileSystem) DeleteEntity(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("erro ao excluir %s: %v", path, err)
	}
	return nil
}

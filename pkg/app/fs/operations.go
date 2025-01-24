package fs

import (
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"fmt"
	"github.com/google/uuid"
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

// --------
//   READ
// --------

// GetFile obtém o caminho de um arquivo com base em seu Id, sua categoria e
// seu usuário, e verifica se o arquivo existe.
func (fs *FileSystem) GetFile(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
	extension string,
) (string, error) {
	// Caminho
	path := fmt.Sprintf(
		"%s/user_%s/categ_%s/file_%s%s",
		fs.Root,
		userId.String(),
		categId.String(),
		fileId.String(),
		extension,
	)

	// Verificar se arquivo existe
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("arquivo %s não existe", path)
	}

	return path, nil
}

// ----------
//   UPDATE
// ----------

// UpdateCategory atualiza o caminho de uma categoria.
func (fs *FileSystem) UpdateCategory(userId uuid.UUID, categId uuid.UUID) error {
	return nil
}

// UpdateFile atualiza o caminho de um arquivo.
func (fs *FileSystem) UpdateFile(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
	extension string,
) error {
	return nil
}

// ----------
//   DELETE
// ----------

// DeleteUser exclui o diretório do usuário com base em seu Id.
func (fs *FileSystem) DeleteUser(userId uuid.UUID) error {
	// Caminho
	path := fmt.Sprintf("%s/user_%s", fs.Root, userId.String())

	// Excluir diretório
	err := DeleteFile(path)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCategory exclui o diretório da categoria com base em seu Id e de seu
// respectivo usuário.
func (fs *FileSystem) DeleteCategory(
	userId uuid.UUID,
	categId uuid.UUID,
) error {
	// Caminho
	path := fmt.Sprintf(
		"%s/user_%s/categ_%s",
		fs.Root,
		userId.String(),
		categId.String(),
	)

	// Excluir diretório
	err := DeleteFile(path)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile exclui um arquivo com base em seu Id e de seu respectivo usuário e
// categoria, com sua extensão.
func (fs *FileSystem) DeleteFile(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
	extension string,
) error {
	// Caminho
	path := fmt.Sprintf(
		"%s/user_%s/categ_%s/file_%s%s",
		fs.Root,
		userId.String(),
		categId.String(),
		fileId.String(),
		extension,
	)

	// Excluir diretório
	err := DeleteFile(path)
	if err != nil {
		return err
	}

	return nil
}

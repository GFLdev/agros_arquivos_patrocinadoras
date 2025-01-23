package fs

import (
	"agros_arquivos_patrocinadoras/filerepo/services/logger"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
)

// ----------
//   CREATE
// ----------

// CreateUser cria um novo usuário no repositório.
func (fs *FileSystem) CreateUser(userId uuid.UUID) error {
	// Caminho
	path := fmt.Sprintf("%s/user_%s", fs.Root, userId.String())

	// Cria o diretório deste usuário
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar pasta do usuário: %v", err)
	}

	return nil
}

// CreateCategory cria uma nova categoria associada a um usuário.
func (fs *FileSystem) CreateCategory(userId uuid.UUID, categId uuid.UUID) error {
	// Caminho
	path := fmt.Sprintf(
		"%s/user_%s/categ_%s",
		fs.Root,
		userId.String(),
		categId.String(),
	)

	// Cria o diretório da nova categoria
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar pasta da categoria: %v", err)
	}

	return nil
}

// CreateFile cria um novo arquivo numa categoria associada a um usuário.
func (fs *FileSystem) CreateFile(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
	extension string,
	content *[]byte,
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

	// Salva o conteúdo deste arquivo em disco
	err := WriteToFile(path, *content, logger.CreateLogger())
	if err != nil {
		return err
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

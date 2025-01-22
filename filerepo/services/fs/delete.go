package fs

import (
	"agros_arquivos_patrocinadoras/filerepo/services/logger"
	"fmt"
	"github.com/google/uuid"
)

func (fs *FS) DeleteUserById(userId uuid.UUID) error {
	user, ok := fs.Users[userId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", userId.String())
	}

	// Excluir diretório
	err := DeleteFiles(user.Path, logger.CreateLogger())
	if err != nil {
		return err
	}

	// Excluir usuário
	delete(fs.Users, userId)
	return StructToFile[FS]("fs/track.json",
		fs,
		logger.CreateLogger(),
	)
}

func (fs *FS) DeleteCategoryById(
	userId uuid.UUID,
	categId uuid.UUID,
) error {
	user, ok := fs.Users[userId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", userId.String())
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return fmt.Errorf("categoria %s não encontrada", categId.String())
	}

	// Excluir diretório
	err := DeleteFiles(categ.Path, logger.CreateLogger())
	if err != nil {
		return err
	}

	// Excluir categoria
	delete(fs.Users[userId].Categories, categId)
	return StructToFile[FS]("fs/track.json",
		fs,
		logger.CreateLogger(),
	)
}

func (fs *FS) DeleteFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) error {
	user, ok := fs.Users[userId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", userId.String())
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return fmt.Errorf("categoria %s não encontrada", categId.String())
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return fmt.Errorf("arquivo %s não encontrado", fileId.String())
	}

	// Excluir arquivo
	err := DeleteFiles(file.Path, logger.CreateLogger())
	if err != nil {
		return err
	}

	// Excluir rastreamento do arquivo
	delete(fs.Users[userId].Categories[categId].Files, fileId)
	return StructToFile[FS]("fs/track.json",
		fs,
		logger.CreateLogger(),
	)
}

package db

import (
	"agros_arquivos_patrocinadoras/logger"
	"fmt"
	"github.com/google/uuid"
)

func (repo *Repo) DeleteUserById(userId uuid.UUID) error {
	user, ok := repo.Users[userId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", userId.String())
	}

	// Excluir diretório
	err := DeleteFiles(user.Path, logger.CreateLogger())
	if err != nil {
		return err
	}

	// Excluir usuário
	delete(repo.Users, userId)
	return StructToFile[Repo]("repo/track.json",
		repo,
		logger.CreateLogger(),
	)
}

func (repo *Repo) DeleteCategoryById(
	userId uuid.UUID,
	categId uuid.UUID,
) error {
	user, ok := repo.Users[userId]
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
	delete(repo.Users[userId].Categories, categId)
	return StructToFile[Repo]("repo/track.json",
		repo,
		logger.CreateLogger(),
	)
}

func (repo *Repo) DeleteFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) error {
	user, ok := repo.Users[userId]
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
	delete(repo.Users[userId].Categories[categId].Files, fileId)
	return StructToFile[Repo]("repo/track.json",
		repo,
		logger.CreateLogger(),
	)
}

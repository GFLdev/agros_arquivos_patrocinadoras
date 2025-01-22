package db

import (
	"agros_arquivos_patrocinadoras/logger"
	"fmt"
	"time"
)

// UpdateUserById atualiza os dados de um usuário com base no ID fornecido.
func (repo *Repo) UpdateUserById(p UpdateUserParams) error {
	// Verificar existência
	user, ok := repo.Users[p.UserId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", p.UserId.String())
	}

	// Atualização
	ts := time.Now().Unix()
	user.Name = p.Name
	user.UpdatedAt = ts
	repo.Users[p.UserId] = user
	repo.UpdatedAt = ts

	// Salvar em disco
	return StructToFile[Repo]("repo/track.json",
		repo,
		logger.CreateLogger(),
	)
}

// UpdateCategoryById atualiza os dados de uma categoria com base no ID fornecido.
func (repo *Repo) UpdateCategoryById(p UpdateCategoryParams) error {
	// Verificar existência
	user, ok := repo.Users[p.UserId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", p.UserId.String())
	}

	categ, ok := user.Categories[p.CategId]
	if !ok {
		return fmt.Errorf("categoria %s não encontrada", p.CategId.String())
	}

	// Atualização
	ts := time.Now().Unix()
	categ.Name = p.Name
	categ.UpdatedAt = ts
	user.Categories[p.CategId] = categ
	user.UpdatedAt = ts
	repo.Users[p.UserId] = user
	repo.UpdatedAt = ts

	// Salvar em disco
	return StructToFile[Repo]("repo/track.json",
		repo,
		logger.CreateLogger(),
	)
}

// UpdateFileById atualiza os dados de um arquivo com base no ID fornecido.
func (repo *Repo) UpdateFileById(p UpdateFileParams) error {
	// Verificar existência
	user, ok := repo.Users[p.UserId]
	if !ok {
		return fmt.Errorf("usuário %s não encontrado", p.UserId.String())
	}

	categ, ok := user.Categories[p.CategId]
	if !ok {
		return fmt.Errorf("categoria %s não encontrada", p.CategId.String())
	}

	file, ok := categ.Files[p.FileId]
	if !ok {
		return fmt.Errorf("arquivo %s não encontrado", p.FileId.String())
	}

	// Atualização
	// Conteúdo
	if len(p.Content) > 0 {
		// Se não houver o Mimetype ou extensão deste novo arquivo
		if p.FileType == "" {
			return fmt.Errorf("mimetype requerido para mudança do arquivo")
		} else {
		}

		if p.Extension == "" {
			return fmt.Errorf("extensão requerida para mudança do arquivo")
		}

		// Escreve em disco. Obs: manter nesta ordem para não alterar a base
		// caso ocorra um erro na gravação do arquivo
		err := WriteToFile(file.Path, p.Content, logger.CreateLogger())
		if err != nil {
			return err
		}

		file.FileType = p.FileType
		file.Extension = p.Extension
	}
	// Nome
	if p.Name != "" {
		file.Name = p.Name
	}

	ts := time.Now().Unix()
	file.UpdatedAt = ts
	categ.Files[p.FileId] = file
	categ.UpdatedAt = ts
	user.Categories[p.CategId] = categ
	user.UpdatedAt = ts
	repo.Users[p.UserId] = user
	repo.UpdatedAt = ts

	// Salvar em disco
	return StructToFile[Repo]("repo/track.json",
		repo,
		logger.CreateLogger(),
	)
}

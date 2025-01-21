package db

import (
	"agros_arquivos_patrocinadoras/logger"
	"fmt"
	"go.uber.org/zap"
	"net/http"
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
	var ext string
	if len(p.Content) > 0 {
		// Se não houver o Mimetype, detectar e definí-lo
		if p.FileType == "" {
			p.FileType = http.DetectContentType(p.Content)
		}
		ext = GetExtension(p.FileType, &p.Content)

		// Escreve em disco. Obs: manter nesta ordem para não alterar a base
		// caso ocorra um erro na gravação do arquivo
		err := WriteToFile(
			file.Path+"."+file.Extension,
			p.Content,
			logger.CreateLogger(),
		)
		if err != nil {
			return err
		}
	}
	// Mimetype
	if p.FileType != "" {
		// Excluir arquivo antigo, caso a extensão tenha mudado
		if ext != file.Extension {
			err := DeleteFiles(file.Path, logger.CreateLogger())
			if err != nil {
				logger.CreateLogger().Warn(
					"Não foi possível excluir o arquivo "+file.Path+"."+file.Extension,
					zap.Error(err),
				)
			}
			file.FileType = p.FileType
			file.Extension = ext
		}
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

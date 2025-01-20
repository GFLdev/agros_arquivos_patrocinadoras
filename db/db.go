package db

import (
	"go.uber.org/zap"
)

func GetFileRepo(logger *zap.Logger) (*Repo, error) {
	// Carregamento do arquivo de rastreamento
	repo, err := FileToStruct[Repo]("repo/track.json", logger)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

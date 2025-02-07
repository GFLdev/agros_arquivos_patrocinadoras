package main

import (
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	types "agros_arquivos_patrocinadoras/pkg/types/config"
	"bytes"
	"crypto/sha256"
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"io"
	"os"
)

func CloseConfigWatcher(ctx *context.Context, watcher *fsnotify.Watcher) {
	if err := watcher.Close(); err != nil {
		ctx.Logger.Error("Erro ao fechar watcher", zap.Error(err))
	}
}

func GetConfigWatcher(ctx *context.Context, restartChan chan bool) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx.Logger.Error("Erro ao criar watcher de "+config.CfgFile, zap.Error(err))
	}

	lastHash, err := calculateFileHash(config.CfgFile)
	if err != nil {
		ctx.Logger.Error("Erro ao calcular hash de "+config.CfgFile, zap.Error(err))
	}

	go func() {
		for {
			bckConfig := ctx.Config
			bckDB := ctx.DB
			select {
			case event := <-watcher.Events:
				// Processa apenas eventos de escrita
				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}

				// Verificar mudanças reais no conteúdo
				currentHash, err := calculateFileHash(config.CfgFile)
				if err != nil {
					ctx.Logger.Error("Erro ao calcular hash de "+event.Name, zap.Error(err))
					continue
				}
				if !hashChanged(lastHash, currentHash) {
					continue
				}

				lastHash = currentHash
				ctx.Logger.Warn("Mudanças detectadas no conteúdo de " + event.Name)

				newConfig, err := config.LoadConfig(ctx.Logger)
				if err != nil {
					ctx.Logger.Error("Erro nas novas configurações. Fallback para o backup", zap.Error(err))
					continue
				}
				ctx.Config = newConfig
				newDB, err := db.GetSqlDB(&newConfig.Database, ctx.Logger)
				if err != nil {
					ctx.Config = bckConfig
					ctx.DB = bckDB
					ctx.Logger.Error("Erro nas novas configurações. Fallback para o backup", zap.Error(err))
					continue
				}
				ctx.Config = newConfig
				ctx.DB = newDB

				// Reiniciar servidor caso os parâmetros para echo tenham alterado
				if serverParamsChanged(bckConfig, newConfig) {
					restartChan <- true
				}
			case err = <-watcher.Errors:
				ctx.Logger.Error("Erro no watcher das configurações", zap.Error(err))
			}
		}
	}()

	return watcher
}

func WatchConfigFile(ctx *context.Context, watcher *fsnotify.Watcher) {
	if err := watcher.Add(config.CfgFile); err != nil {
		ctx.Logger.Error("Erro ao adicionar "+config.CfgFile+" ao watcher", zap.Error(err))
	}
}

func calculateFileHash(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	return hash[:], nil
}

func hashChanged(oldHash []byte, newHash []byte) bool {
	return !bytes.Equal(oldHash, newHash)
}

func serverParamsChanged(oldConfig, newConfig *types.Config) bool {
	return newConfig.Port != oldConfig.Port ||
		newConfig.Environment != oldConfig.Environment ||
		(newConfig.Environment == "development" &&
			(newConfig.CertFile != oldConfig.CertFile ||
				newConfig.KeyFile != oldConfig.KeyFile))

}

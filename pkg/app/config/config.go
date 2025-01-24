package config

import (
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
)

// LoadConfig carrega as configurações do servidor.
func LoadConfig(logr *zap.Logger) *config.Config {
	logr.Info("Carregando arquivo de configurações")

	// Abertura do arquivo de configuração
	file, err := os.Open("cfg.json")
	if err != nil {
		logr.Fatal("Não foi possível abrir arquivo de configuração",
			zap.Error(err),
		)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logr.Fatal("Não foi possível fechar arquivo de configuração",
				zap.Error(err),
			)
		}
	}(file)

	// Leitura do arquivo
	jsonPayload, err := io.ReadAll(file)
	if err != nil {
		logr.Fatal("Não foi possível ler arquivo de configuração",
			zap.Error(err),
		)
	}

	// Desagrupamento do JSON
	cfg := &config.Config{}
	err = json.Unmarshal(jsonPayload, cfg)
	if err != nil {
		logr.Fatal("Não foi possível desagrupar dados de configuração",
			zap.Error(err),
		)
	}

	// Logging
	if cfg.Environment == "production" {
		logr.Info("Configurando servidor de produção")
	} else if cfg.Environment == "development" {
		logr.Info("Configurando servidor de desenvolvimento")
	} else {
		logr.Warn("Ambiente não definido. Fallback para development")
		cfg.Environment = "development"
	}

	return cfg
}

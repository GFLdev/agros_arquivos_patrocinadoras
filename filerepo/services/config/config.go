package config

import (
	"agros_arquivos_patrocinadoras/filerepo/db"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
)

type Config struct {
	Environment string      `json:"environment" validate:"required"`
	Origins     []string    `json:"origins" validate:"required"`
	Port        int         `json:"port" validate:"required"`
	Database    db.Database `json:"database" validate:"required"`
	JwtSecret   string      `json:"jwt_secret" validate:"required"`
	JwtExpires  int         `json:"jwt_expires" validate:"required"`
	CertFile    string      `json:"cert_file"`
	KeyFile     string      `json:"key_file"`
}

// LoadConfig carrega as configurações do servidor.
func LoadConfig(logr *zap.Logger) *Config {
	logr.Info("Carregando arquivo de configurações")

	// Abertura do arquivo de configuração
	file, err := os.Open("config.json")
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
	config := &Config{}
	err = json.Unmarshal(jsonPayload, config)
	if err != nil {
		logr.Fatal("Não foi possível desagrupar dados de configuração",
			zap.Error(err),
		)
	}

	// Logging
	if config.Environment == "production" {
		logr.Info("Configurando servidor de produção")
	} else if config.Environment == "development" {
		logr.Info("Configurando servidor de desenvolvimento")
	} else {
		logr.Warn("Ambiente não definido. Fallback para development")
		config.Environment = "development"
	}

	return config
}

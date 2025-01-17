package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os"
)

// LoadConfig carrega as configurações do servidor.
func LoadConfig(ctx *AppContext) *Config {
	ctx.Logger.Info("Carregando arquivo de configurações")

	// Abertura do arquivo de configuração
	file, err := os.Open("config.json")
	if err != nil {
		ctx.Logger.Fatal("Não foi possível abrir arquivo de configuração",
			zap.Error(err),
		)
	}
	defer file.Close()

	// Leitura do arquivo
	jsonPayload, err := io.ReadAll(file)
	if err != nil {
		ctx.Logger.Fatal("Não foi possível ler arquivo de configuração",
			zap.Error(err),
		)
	}

	// Desagrupamento do JSON
	config := &Config{}
	err = json.Unmarshal(jsonPayload, config)
	if err != nil {
		ctx.Logger.Fatal("Não foi possível desagrupar dados de configuração",
			zap.Error(err),
		)
	}

	// Logging
	if config.Environment == "production" {
		ctx.Logger.Info("Configurando servidor de produção")
	} else if config.Environment == "development" {
		ctx.Logger.Info("Configurando servidor de desenvolvimento")
	} else {
		ctx.Logger.Warn("Ambiente não definido. Fallback para development")
		config.Environment = "development"
	}

	return config
}

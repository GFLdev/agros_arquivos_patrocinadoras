// Package config fornece funcionalidades para carregar e gerenciar configurações
// da aplicação a partir de um arquivo JSON. Ele inclui métodos para lidar com a
// leitura, desagrupamento e validação de dados de configuração, garantindo que o
// ambiente da aplicação seja configurado corretamente.
//
// Este pacote foi projetado para usar o logger zap para fornecer informações
// detalhadas sobre falhas e erros relacionados às configurações.
package config

import (
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
)

// CfgFile é o nome padrão do arquivo de configuração JSON da aplicação.
const CfgFile = "config.json"

// LoadConfig lê e carrega as configurações da aplicação a partir de um arquivo
// JSON padrão.
//
// Parâmetros:
//   - logr: zap.Logger da aplicação para logging.
//
// Retorno:
//   - *config.Config: ponteiro para a estrutura de configuração carregada.
//   - error: erro caso tenha acontecido algum erro no carregamento.
func LoadConfig(logr *zap.Logger) (*config.Config, error) {
	logr.Info("Carregando arquivo de configurações")

	// Abertura do arquivo de configuração
	file, err := os.Open(CfgFile)
	if err != nil {
		return nil, fmt.Errorf("não foi possível abrir arquivo de configuração: %w", err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			logr.Error("Não foi possível fechar arquivo de configuração", zap.Error(err))
		}
	}(file)

	// Leitura do arquivo e desagrupamento do JSON
	jsonPayload, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("não foi possível ler arquivo de configuração: %w", err)
	}
	cfg := &config.Config{}
	if err = json.Unmarshal(jsonPayload, cfg); err != nil {
		return nil, fmt.Errorf("não foi possível desagrupar dados de configuração: %w", err)
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

	return cfg, nil
}

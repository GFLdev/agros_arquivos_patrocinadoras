package handlers

import (
	"agros_arquivos_patrocinadoras/db"
	"go.uber.org/zap"
	"sync"
)

type AppContext struct {
	Logger *zap.Logger
	Config *Config
	Repo   struct {
		*db.Repo
		*sync.Mutex
	}
}

type Config struct {
	Environment string   `json:"environment" validate:"required"`
	Origins     []string `json:"origins" validate:"required"`
	Port        int      `json:"port" validate:"required"`
	CertFile    string   `json:"cert_file"`
	KeyFile     string   `json:"key_file"`
}

type ErrorRes struct {
	Message string `json:"message" validate:"required"`
	Error   error  `json:"error" validate:"required"`
}

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
	User          string `json:"user" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
}

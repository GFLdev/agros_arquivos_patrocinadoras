package main

import (
	"go.uber.org/zap"
)

type AppContext struct {
	Logger *zap.Logger
	Config *Config
}

type Config struct {
	Environment string   `json:"environment" validate:"required"`
	Origins     []string `json:"origins" validate:"required"`
	Port        int      `json:"port" validate:"required"`
	LogFile     string   `json:"log_file" validate:"required"`
	CertFile    string   `json:"cert_file"`
	KeyFile     string   `json:"key_file"`
}

type ErrorResponse struct {
	Message string `json:"message" validate:"required"`
	Error   error  `json:"error" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User          string `json:"user" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
}

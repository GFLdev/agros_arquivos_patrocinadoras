package handlers

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// BodyUnmarshall lê e desagrupa o body de uma requisição.
func BodyUnmarshall[T any](c echo.Context) (*T, error) {
	// Obtenção do logger
	logr := context.GetContext(c).Logger

	// Leitura do body como JSON
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logr.Error("Erro ao ler o corpo da requisição", zap.Error(err))
		return nil, err
	}

	// Desagrupamento e retorno
	result := new(T)
	err = json.Unmarshal(payload, result)
	if err != nil {
		logr.Error("Erro ao desagrupar o corpo da requisição", zap.Error(err))
		return nil, err
	}
	return result, nil
}

// LogHTTPDetails registra os detalhes de uma requisição HTTP, incluindo o IP
// do cliente, o método e o caminho, com campos personalizados adicionais.
func LogHTTPDetails(c echo.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	// Obtenção do logger
	logr := context.GetContext(c).Logger

	// Informações da requisição/resposta
	baseFields := []zap.Field{
		zap.String("method", c.Request().Method),
		zap.String("path", c.Path()),
		zap.String("client_ip", c.RealIP()),
	}
	allFields := append(baseFields, fields...)

	// Logger com o nível escolhido
	switch level {
	case zapcore.DebugLevel:
		logr.Debug(msg, allFields...)
	case zapcore.InfoLevel:
		logr.Info(msg, allFields...)
	case zapcore.WarnLevel:
		logr.Warn(msg, allFields...)
	case zapcore.ErrorLevel:
		logr.Error(msg, allFields...)
	case zapcore.FatalLevel:
		logr.Fatal(msg, allFields...)
	default:
		logr.Info(msg, allFields...)
	}
}

// HashPassword gera o hash da senha usando bcrypt com um custo padrão
func HashPassword(ctx *context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.Logger.Error("Erro ao gerar hash da senha.", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

func ParseEntityUUID(c echo.Context, entityType fs.EntityType) (uuid.UUID, error) {
	var param string
	switch entityType {
	case fs.User:
		param = c.Param("userId")
	case fs.Category:
		param = c.Param("categId")
	case fs.File:
		param = c.Param("fileId")
	default:
		return uuid.Nil, fmt.Errorf("entidade %d não suportada", entityType)
	}

	id, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, fmt.Errorf("uuid inválido: %v", err)
	}
	return id, nil
}

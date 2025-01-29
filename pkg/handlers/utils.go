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

// BodyUnmarshall realiza o desagrupamento (unmarshal) do corpo da requisição
// para uma estrutura genérica.
//
// Parâmetros:
//   - c: contexto da requisição HTTP (do pacote echo).
//
// Retornos:
//   - *T: ponteiro para a estrutura desagrupada.
//   - error: erro, caso ocorra ao ler ou desagrupar o corpo da requisição.
func BodyUnmarshall[T any](c echo.Context) (*T, error) {
	// Obtenção do logger e leitura do body como JSON
	logr := context.GetContext(c).Logger
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

// LogHTTPDetails registra os detalhes de uma solicitação HTTP no logger
// associado ao contexto da aplicação, com o nível de log especificado.
//
// Parâmetros:
//   - c: contexto da requisição HTTP (do pacote echo).
//   - level: nível do log (Debug, Info, Warn, Error, Fatal).
//   - msg: mensagem a ser registrada.
//   - fields: campos adicionais para detalhar o log.
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

// HashPassword gera um hash seguro para a senha fornecida utilizando o bcrypt.
//
// Parâmetros:
//   - ctx: contexto da aplicação, contendo o logger para registrar possíveis
//     erros.
//   - password: string contendo a senha a ser hasheada.
//
// Retornos:
//   - string: hash gerado a partir da senha fornecida.
//   - error: erro, caso ocorra falha ao gerar o hash.
func HashPassword(ctx *context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.Logger.Error("Erro ao gerar hash da senha.", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

// ParseEntityUUID realiza o parse do UUID baseado no tipo de entidade e no
// parâmetro correspondente presente na requisição HTTP.
//
// Parâmetros:
//   - c: contexto da requisição HTTP (do pacote echo).
//   - entityType: tipo da entidade que define qual parâmetro UUID será lido
//     (fs.User, fs.Category, fs.File).
//
// Retornos:
//   - uuid.UUID: o UUID extraído e parseado do parâmetro.
//   - error: erro, caso o UUID não seja encontrado ou inválido.
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
		return uuid.Nil, fmt.Errorf("uuid inválido: %w", err)
	}
	return id, nil
}

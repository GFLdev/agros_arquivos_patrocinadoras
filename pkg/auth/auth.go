// Package auth fornece funcionalidades relacionadas à autenticação e geração de
// tokens JWT, incluindo a definição de claims personalizadas e métodos para
// criar tokens seguros usando o algoritmo de assinatura HS256.
package auth

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

// CustomClaims define uma estrutura personalizada para os claims de um token
// JWT.
type CustomClaims struct {
	// Id representa o identificador único para os claims, usado para
	// identificar um usuário.
	Id uuid.UUID `json:"id"`
	// Name representa o nome do usuário ou entidade associado aos claims.
	Name string `json:"name"`
	// Admin indica se o usuário possui privilégios administrativos.
	Admin bool `json:"admin"`
	jwt.RegisteredClaims
}

// GenerateToken cria um token JWT com claims personalizados, utilizando o
// algoritmo HS256.
//
// Parâmetros:
//   - c: contexto das requisições HTTP que contém informações
//     do request atual.
//   - userId: identificador único do usuário.
//   - userName: nome do usuário.
//
// Retornos:
//   - string: token JWT gerado.
//   - error: erro caso ocorra algum problema durante a geração do token.
func GenerateToken(c echo.Context, userId uuid.UUID, userName string) (string, error) {
	ctx := context.GetContext(c)

	// JWT Claims
	duration := time.Duration(ctx.Config.JwtExpires)
	claims := CustomClaims{
		Id:    userId,
		Name:  userName,
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(duration * time.Minute),
			),
		},
	}

	// Declaração do token com algoritmo HS256 e geração do token criptografado
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(ctx.Config.JwtSecret))
	if err != nil {
		ctx.Logger.Error("Erro ao gerar JWT.", zap.Error(err))
		return "", err
	}
	return t, nil
}

// Package auth fornece funcionalidades relacionadas à autenticação e geração de
// tokens JWT, incluindo a definição de claims personalizadas e métodos para
// criar tokens seguros usando o algoritmo de assinatura HS256.
package auth

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

// ClaimsData armazena informações relacionadas a um usuário para os Claims de
// seu token.
type ClaimsData struct {
	// Id representa o identificador único para os claims, usado para
	// identificar um usuário.
	Id uuid.UUID `json:"id"`
	// Name representa o nome de apresentação do usuário.
	Name string `json:"name"`
}

// CustomClaims define uma estrutura personalizada para os claims de um token
// JWT.
type CustomClaims struct {
	ClaimsData
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
func GenerateToken(c echo.Context, data ClaimsData, expiresAt time.Time) (string, error) {
	ctx := context.GetContext(c)

	// JWT Claims
	claims := CustomClaims{
		ClaimsData: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
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

func GetClaims(c echo.Context) (*CustomClaims, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("erro ao obter token")
	}

	claims, ok := user.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("erro ao obter claims")
	}
	return claims, nil
}

func AuthenticateAdmin(c echo.Context) bool {
	ctx := context.GetContext(c)
	claims, err := GetClaims(c)
	if err != nil || claims.Id != ctx.AdminId {
		return false
	}
	return true
}

func AuthenticateUser(c echo.Context, userId uuid.UUID) bool {
	ctx := context.GetContext(c)
	claims, err := GetClaims(c)
	if err != nil {
		return false
	}

	if claims.Id != userId && claims.Id != ctx.AdminId {
		return false
	}

	return true
}

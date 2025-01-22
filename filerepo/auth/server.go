package auth

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

type CustomClaims struct {
	UserId   uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Admin    bool      `json:"admin"`
	jwt.RegisteredClaims
}

// GenerateToken gera token JWT.
func GenerateToken(l *utils.LoginReq, c echo.Context) (string, error) {
	ctx := services.GetContext(c)

	// FIXME: Criar lógica de autenticação do usuário
	userId := uuid.New()

	// JWT Claims
	duration := time.Duration(ctx.Config.JwtExpires)
	claims := CustomClaims{
		UserId:   userId,
		Username: l.Username,
		Admin:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(duration * time.Minute),
			),
		},
	}

	// Declaração do token com algoritmo HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Gerar token criptografado
	t, err := token.SignedString([]byte(ctx.Config.JwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

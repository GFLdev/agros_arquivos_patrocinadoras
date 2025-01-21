package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

const (
	accessTokenCookieName = "access-token"
)

type Claims struct {
	Name string `json:"name"`
	jwt.Claims
}

// GetJWTSecret retorna o segredo do JWT.
func GetJWTSecret() string {
	return viper.GetString("JWT_SECRET")
}

// GetJWTExpirationTime retorna a duração do JWT.
func GetJWTExpirationTime() time.Duration {
	return viper.GetDuration("JWT_EXPIRATION_TIME")
}

// GenerateToken gera token JWT.
func GenerateToken(
	l *LoginReq,
	c echo.Context,
	exp time.Time,
) (string, error) {
	// JWT Claims
	claims := jwt.MapClaims{
		"username":  l.Username,
		"password":  l.Password,
		"expiresAt": exp.Unix(),
	}

	// Declaração do token com algoritmo HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// setTokenCookie cria cookie que armazena o token JWT.
func setTokenCookie(c echo.Context, token string, exp time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = accessTokenCookieName
	cookie.Value = token
	cookie.Expires = exp
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// setUserCookie cria cookie que armazena o token JWT.
func setUserCookie(c echo.Context, l *LoginReq, exp time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = l.Username
	cookie.Expires = exp
	cookie.Path = "/"

	c.SetCookie(cookie)
}

func Wrapper() {
	//expiration := time.Now().Add(GetJWTExpirationTime())
}

func IsAuthenticated(c echo.Context) bool {
	return true
}

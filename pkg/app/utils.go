package app

import (
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

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

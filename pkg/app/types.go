package app

import (
	"github.com/google/uuid"
)

// LoginParams define os parâmetros para o login de um usuário.
type LoginParams struct {
	// Username especifica o nome do usuário.
	Username string
	// Password especifica a senha do usuário.
	Password string
}

// LoginData define os dados de resposta da criação de um usuário.
type LoginData struct {
	// UserId especifica o identificador do usuário.
	UserId uuid.UUID
	// Name especificar o nome de apresentação do usuário
	Name string
}

// UserData define os parâmetros para a criação de um usuário.
type UserData struct {
	// Username especifica o nome do usuário.
	Username string
	// Name especificar o nome de apresentação do usuário
	Name string
	// Password especifica a senha do usuário.
	Password string
}

// CategData define os parâmetros para a criação de uma categoria.
type CategData struct {
	// UserId especifica o identificador do usuário associado à categoria.
	UserId uuid.UUID
	// Name especifica o nome da categoria.
	Name string
}

// FileData define os parâmetros para a criação de um arquivo.
type FileData struct {
	// CategId especifica o identificador da categoria associada ao arquivo.
	CategId uuid.UUID
	// Name especifica o nome do arquivo.
	Name string
	// Extension especifica a extensão do arquivo.
	Extension string
	// Mimetype especifica o tipo MIME do arquivo.
	Mimetype string
	// Content contém o conteúdo do arquivo.
	Content *[]byte
}

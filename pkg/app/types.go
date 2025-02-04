package app

import (
	"github.com/google/uuid"
)

// UserParams define os parâmetros para a criação de um usuário.
type UserParams struct {
	// Name especifica o nome do usuário.
	Name string
	// Password especifica a senha do usuário.
	Password string
}

// CategParams define os parâmetros para a criação de uma categoria.
type CategParams struct {
	// UserId especifica o identificador do usuário associado à categoria.
	UserId uuid.UUID
	// Name especifica o nome da categoria.
	Name string
}

// FileParams define os parâmetros para a criação de um arquivo.
type FileParams struct {
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

// LoginCompare define os parâmetros necessários para realizar a comparação de
// login de um usuário.
type LoginCompare struct {
	// UserId especifica o identificador do usuário.
	UserId uuid.UUID
	// Hash representa a senha criptografada associada ao usuário.
	Hash string
}

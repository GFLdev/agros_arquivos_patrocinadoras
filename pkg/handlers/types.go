package handlers

import "github.com/google/uuid"

// LoginReq representa os dados necessários para autenticação de um usuário.
type LoginReq struct {
	// Username especifica o nome de usuário para autenticação.
	Username string `json:"username" validate:"required"`
	// Password especifica a senha para autenticação.
	Password string `json:"password" validate:"required"`
}

// CreateUserReq representa os dados necessários para criar um novo usuário.
type CreateUserReq struct {
	// Name especifica o nome do novo usuário.
	Name string `json:"name" validate:"required"`
	// Password especifica a senha do novo usuário.
	Password string `json:"password" validate:"required"`
}

// CreateCategoryReq representa os dados necessários para criar uma nova categoria.
type CreateCategoryReq struct {
	// Name especifica o nome da nova categoria.
	Name string `json:"name" validate:"required"`
}

// CreateFileReq representa os dados necessários para criar um novo arquivo.
type CreateFileReq struct {
	// Name especifica o nome do novo arquivo.
	Name string `json:"name" validate:"required"`
	// Extension especifica a extensão do novo arquivo.
	Extension string `json:"extension" validate:"required"`
	// Mimetype especifica o tipo MIME do novo arquivo.
	Mimetype string `json:"mimetype" validate:"required"`
	// Content especifica o conteúdo do novo arquivo.
	Content []byte `json:"content" validate:"required"`
}

// UpdateUserReq representa os dados necessários para atualizar um usuário.
type UpdateUserReq struct {
	// Name especifica o novo nome do usuário.
	Name string `json:"name" validate:"required"`
	// Password especifica a nova senha do usuário.
	Password string `json:"password" validate:"required"`
}

// UpdateCategoryReq representa os dados necessários para atualizar uma categoria.
type UpdateCategoryReq struct {
	// UserId especifica o ID do usuário proprietário da categoria.
	UserId uuid.UUID `json:"userId" validate:"required"`
	// Name especifica o novo nome da categoria.
	Name string `json:"name" validate:"required"`
}

// UpdateFileReq representa os dados necessários para atualizar um arquivo.
type UpdateFileReq struct {
	// CategId especifica o ID da categoria associada ao arquivo.
	CategId uuid.UUID `json:"categId" validate:"required"`
	// Name especifica o novo nome do arquivo.
	Name string `json:"name" validate:"required"`
	// Extension especifica a nova extensão do arquivo.
	Extension string `json:"extension" validate:"required"`
	// Mimetype especifica o novo tipo MIME do arquivo.
	Mimetype string `json:"mimetype" validate:"required"`
	// Content especifica o novo conteúdo do arquivo.
	Content []byte `json:"content" validate:"required"`
}

// ErrorRes representa a estrutura de uma resposta de erro.
type ErrorRes struct {
	// Message especifica a mensagem descritiva do erro.
	Message string `json:"message" validate:"required"`
	// Error especifica os detalhes técnicos do erro.
	Error string `json:"error" validate:"required"`
}

// GenericRes representa a estrutura de uma resposta genérica.
type GenericRes struct {
	// Message especifica a mensagem descritiva da resposta.
	Message string `json:"message" validate:"required"`
}

// LoginRes representa a resposta para uma requisição de login.
type LoginRes struct {
	// User especifica o nome do usuário autenticado.
	User string `json:"user" validate:"required"`
	// Authenticated indica se o usuário foi autenticado com sucesso.
	Authenticated bool `json:"authenticated" validate:"required"`
}

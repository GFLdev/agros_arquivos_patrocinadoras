package handlers

import "github.com/google/uuid"

// EntityType define os tipos de entidades possíveis no sistema.
type EntityType int

const (
	// User representa um tipo de entidade para usuários.
	User EntityType = iota
	// Category representa um tipo de entidade para categorias.
	Category
	// File representa um tipo de entidade para arquivos.
	File
)

// LoginReq representa os dados necessários para autenticação de um usuário.
type LoginReq struct {
	// Username especifica o nome de usuário para autenticação.
	Username string `json:"username" validate:"required"`
	// Password especifica a senha para autenticação.
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
	Token   string      `json:"token"`
	Message HTTPMessage `json:"message"`
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Admin   bool        `json:"admin"`
}

// CreateUserReq representa os dados necessários para criar um novo usuário.
type CreateUserReq struct {
	// Username especifica o nome de usuário, do novo usuário.
	Username string `json:"username" validate:"required"`
	// Name especifica o nome de apresentação do novo usuário.
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
	// Username especifica o novo nome de usuário.
	Username string `json:"username"`
	// Name especifica o novo nome de apresentação do usuário.
	Name string `json:"name"`
	// Password especifica a nova senha do usuário.
	Password string `json:"password"`
}

// UpdateCategoryReq representa os dados necessários para atualizar uma categoria.
type UpdateCategoryReq struct {
	// UserId especifica o ID do usuário proprietário da categoria.
	UserId string `json:"user_id"`
	// Name especifica o novo nome da categoria.
	Name string `json:"name"`
}

// UpdateFileReq representa os dados necessários para atualizar um arquivo.
type UpdateFileReq struct {
	// CategId especifica o ID da categoria associada ao arquivo.
	CategId string `json:"categ_id"`
	// Name especifica o novo nome do arquivo.
	Name string `json:"name"`
	// Extension especifica a nova extensão do arquivo.
	Extension string `json:"extension"`
	// Mimetype especifica o novo tipo MIME do arquivo.
	Mimetype string `json:"mimetype"`
	// Content especifica o novo conteúdo do arquivo.
	Content []byte `json:"content"`
}

// CreateResponse representa a resposta retornada após uma operação de criação
// bem sucedida.
type CreateResponse struct {
	// Id é o identificador único da entidade.
	Id uuid.UUID `json:"id"`
	// Message é a descrição de retorno da operação.
	Message HTTPMessage `json:"message"`
}

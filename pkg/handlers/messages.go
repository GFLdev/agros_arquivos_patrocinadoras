package handlers

type HTTPMessage string

// Mensagens relacionadas ao usuário.
const (
	InvalidUserIdMessage HTTPMessage = "Id de usuário inválido."
	CreatedUserMessage   HTTPMessage = "Usuário criado com sucesso."
	UserNotFoundMessage  HTTPMessage = "Usuário não encontrado."
	UsersNotFoundMessage HTTPMessage = "Nenhum usuário foi encontrado."
	UpdatedUserMessage   HTTPMessage = "Usuário atualizado com sucesso."
	DeletedUserMessage   HTTPMessage = "Usuário excluído com sucesso."
)

// Mensagens relacionadas à categoria.
const (
	InvalidCategoryIdMessage  HTTPMessage = "Id de categoria inválido."
	CreatedCategoryMessage    HTTPMessage = "Categoria criada com sucesso."
	CategoryNotFoundMessage   HTTPMessage = "Categoria não encontrada."
	CategoriesNotFoundMessage HTTPMessage = "Nenhuma categoria foi encontrada."
	UpdatedCategoryMessage    HTTPMessage = "Categoria atualizada com sucesso."
	DeletedCategoryMessage    HTTPMessage = "Categoria excluída com sucesso."
)

// Mensagens relacionadas ao arquivo.
const (
	InvalidFileIdMessage HTTPMessage = "Id de arquivo inválido."
	CreatedFileMessage   HTTPMessage = "Arquivo criado com sucesso."
	FileNotFoundMessage  HTTPMessage = "Arquivo não encontrado."
	FilesNotFoundMessage HTTPMessage = "Nenhum arquivo foi encontrado."
	UpdatedFileMessage   HTTPMessage = "Arquivo atualizado com sucesso."
	DeletedFileMessage   HTTPMessage = "Arquivo excluído com sucesso."
)

// Mensagens gerais.
const (
	BadRequestMessage          HTTPMessage = "Falha na requisição. Verifique os dados e tente novamente."
	InternalServerErrorMessage HTTPMessage = "Erro interno no sistema. Tente novamente."
	UnauthorizedMessage        HTTPMessage = "Acesso negado. Verifique suas credenciais."
)

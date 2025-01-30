package app

import (
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"database/sql"
	"github.com/google/uuid"
)

// RollbackErrors armazena erros ocorridos durante o rollback de transações.
type RollbackErrors struct {
	// DB especifica o erro relacionado ao banco de dados.
	DB error
	// FS especifica o erro relacionado ao sistema de arquivos.
	FS error
}

// CreateRollbackData contém os dados necessários para realizar o rollback de
// uma criação.
type CreateRollbackData struct {
	// Tx representa a transação do banco de dados.
	Tx *sql.Tx
	// Path especifica o caminho no sistema de arquivos.
	Path string
	// RollbackErrors armazena os erros associados ao rollback.
	*RollbackErrors
}

// UpdateRollbackData contém os dados necessários para realizar o rollback de
// uma atualização.
type UpdateRollbackData struct {
	// Tx representa a transação do banco de dados.
	Tx *sql.Tx
	// OldPath especifica o caminho antigo no sistema de arquivos.
	OldPath string
	// NewPath especifica o novo caminho no sistema de arquivos.
	NewPath string
	// RollbackErrors armazena os erros associados ao rollback.
	*RollbackErrors
}

// DeleteRollbackData contém os dados necessários para realizar o rollback de
// uma exclusão.
type DeleteRollbackData struct {
	// Tx representa a transação do banco de dados.
	Tx *sql.Tx
	// Path especifica o caminho no sistema de arquivos.
	Path string
	// Type indica o tipo de entidade sendo excluída.
	Type fs.EntityType
	// Content armazena o conteúdo do arquivo excluído.
	Content *[]byte
	// RollbackErrors armazena os erros associados ao rollback.
	*RollbackErrors
}

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
	// UserId especifica o identificador do usuário associado ao arquivo.
	UserId uuid.UUID
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

// UserUpdate define os parâmetros para a atualização de um usuário.
type UserUpdate struct {
	// UserId especifica o identificador do usuário a ser atualizado.
	UserId uuid.UUID
	// OldName especifica o novo nome do usuário.
	OldName string
	// UserParams contém os parâmetros a mudar do usuário.
	UserParams
}

// CategUpdate define os parâmetros para a atualização de uma categoria.
type CategUpdate struct {
	// CategId especifica o identificador da categoria a ser atualizada.
	CategId uuid.UUID
	// OldUserId especifica o novo identificador do usuário associado.
	OldUserId uuid.UUID
	// OldName especifica o novo nome da categoria.
	OldName string
	// CategParams contém os parâmetros a mudar da categoria.
	CategParams
}

// FileUpdate define os parâmetros para a atualização de um arquivo.
type FileUpdate struct {
	// FileId especifica o identificador do arquivo a ser atualizado.
	FileId uuid.UUID
	// OldCategId especifica o novo identificador da categoria associada.
	OldCategId uuid.UUID
	// OldName especifica o novo nome do arquivo.
	OldName string
	// OldExtension especifica a nova extensão do arquivo.
	OldExtension string
	// OldMimetype especifica o novo tipo MIME do arquivo.
	OldMimetype string
	// FileParams contém os parâmetros a mudar do arquivo.
	FileParams
}

// UserDelete define os parâmetros para a exclusão de um usuário.
type UserDelete struct {
	// UserId especifica o identificador do usuário a ser excluído.
	UserId uuid.UUID
}

// CategDelete define os parâmetros para a exclusão de uma categoria.
type CategDelete struct {
	// UserDelete contém os parâmetros adicionais da exclusão de usuário.
	UserDelete
	// CategId especifica o identificador da categoria a ser excluída.
	CategId uuid.UUID
}

// FileDelete define os parâmetros para a exclusão de um arquivo.
type FileDelete struct {
	// CategDelete contém os parâmetros adicionais da exclusão de categoria.
	CategDelete
	// FileId especifica o identificador do arquivo a ser excluído.
	FileId uuid.UUID
	// Extension especifica a extensão do arquivo a ser excluído.
	Extension string
}

// LoginCompare define os parâmetros necessários para realizar a comparação de
// login de um usuário.
type LoginCompare struct {
	// UserId especifica o identificador do usuário.
	UserId uuid.UUID
	// Hash representa a senha criptografada associada ao usuário.
	Hash string
}

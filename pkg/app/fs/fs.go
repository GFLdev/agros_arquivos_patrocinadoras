// Package fs fornece definições para gerenciamento de entidades de um sistema
// de arquivos, como usuários, categorias, arquivos e anexos. Este pacote inclui
// estruturas que representam dados comuns às entidades, bem como informações
// específicas sobre arquivos, como extensão, tipo MIME e caminho no sistema de
// arquivos.
package fs

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

// FileSystem representa o sistema de arquivos, incluindo o diretório raiz.
type FileSystem struct {
	// Root especifica o diretório raiz do sistema de arquivos.
	Root string `json:"root"`
}

// EntityData contém informações comuns sobre uma entidade, como id, nome e
// data de atualização.
type EntityData struct {
	// Id é o identificador único da entidade.
	Id uuid.UUID `json:"id" validate:"required"`
	// Name é o nome da entidade.
	Name string `json:"name" validate:"required"`
	// UpdatedAt é o timestamp da última atualização da entidade.
	UpdatedAt int64 `json:"updated_at" validate:"required"`
}

// FileData contém informações sobre um arquivo, incluindo dados gerais da
// entidade, extensão do arquivo e tipo MIME.
type FileData struct {
	// EntityData contém informações gerais sobre a entidade, como ID, nome e
	// data de atualização.
	EntityData
	// Extension é a extensão do arquivo, como ".jpg", ".txt", etc.
	Extension string `json:"extension" validate:"required"`
	// Mimetype é o tipo MIME do arquivo, como "image/jpeg", "text/plain", etc.
	Mimetype string `json:"mimetype" validate:"required"`
}

// AttachmentData contém informações sobre um anexo, incluindo os dados do
// arquivo e o caminho onde o arquivo está armazenado.
type AttachmentData struct {
	// FileData contém as informações gerais sobre o arquivo, como identificador,
	// nome, extensão, tipo MIME e data de atualização.
	FileData
	// Path é o caminho do arquivo do anexo no sistema de arquivos.
	Path string `json:"path" validate:"required"`
}

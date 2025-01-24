package fs

import "github.com/google/uuid"

type EntityType int

const (
	User EntityType = iota
	Category
	File
)

type FileSystem struct {
	Root string `json:"root"`
}

// EntityData representa a estrutura básica para realizar consultas.
type EntityData struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	UpdatedAt int64     `json:"updated_at" validate:"required"`
}

// FileData estende EntityData para incluir informações específicas de
// arquivos.
type FileData struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	FileType  string    `json:"file_type" validate:"required"`
	Extension string    `json:"extension" validate:"required"`
	UpdatedAt int64     `json:"updated_at" validate:"required"`
}

// AttachmentData estende FileData para incluir informações específicas de
// arquivos para download.
type AttachmentData struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	FileType  string    `json:"file_type" validate:"required"`
	Path      string    `json:"path" validate:"required"`
	Extension string    `json:"extension" validate:"required"`
	UpdatedAt int64     `json:"updated_at" validate:"required"`
}

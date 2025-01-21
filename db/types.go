package db

import (
	"github.com/google/uuid"
)

// ---------------
//   Repositório
// ---------------

// File representa a estrutura de um arquivo no sistema.
type File struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	FileType  string    `json:"file_type" validate:"required"`
	Path      string    `json:"path" validate:"required"`
	Extension string    `json:"extension" validate:"required"`
	UpdatedAt int64     `json:"updated_at" validate:"required"`
}

// Category representa uma coleção de arquivos organizados sob um nome e caminho
// específicos.
type Category struct {
	Id        uuid.UUID          `json:"id" validate:"required"`
	Name      string             `json:"name" validate:"required"`
	Path      string             `json:"path" validate:"required"`
	Files     map[uuid.UUID]File `json:"files" validate:"required"`
	UpdatedAt int64              `json:"updated_at" validate:"required"`
}

// User representa um usuário no sistema, contendo as categorias associadas.
type User struct {
	Id         uuid.UUID              `json:"id" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	Path       string                 `json:"path" validate:"required"`
	Categories map[uuid.UUID]Category `json:"categories" validate:"required"`
	UpdatedAt  int64                  `json:"updated_at" validate:"required"`
}

// Repo representa o repositório contendo todos os usuários e metadados.
type Repo struct {
	Users     map[uuid.UUID]User `json:"users" validate:"required"`
	UpdatedAt int64              `json:"updated_at" validate:"required"`
}

// ----------
//   CREATE
// ----------

// CreateUserParams representa os dados necessários para criar um novo usuário.
type CreateUserParams struct {
	Name string `json:"name" validate:"required"`
}

// CreateCategoryParams representa os dados necessários para criar uma nova
// categoria.
type CreateCategoryParams struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Name   string    `json:"name" validate:"required"`
}

// CreateFileParams representa os dados necessários para criar um novo arquivo.
type CreateFileParams struct {
	UserId   uuid.UUID `json:"user_id" validate:"required"`
	CategId  uuid.UUID `json:"categ_id" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	FileType string    `json:"file_type" validate:"required"`
	Content  []byte    `json:"content" validate:"required"`
}

// --------
//   READ
// --------

// QueryData representa a estrutura básica para realizar consultas.
type QueryData struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	UpdatedAt int64     `json:"updated_at" validate:"required"`
}

// FileAttachment estende QueryData para incluir informações específicas de
// arquivos.
type FileAttachment struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	FileType  string    `json:"file_type" validate:"required"`
	Path      string    `json:"path" validate:"required"`
	Extension string    `json:"extension" validate:"required"`
	UpdatedAt int64     `json:"updated_at" validate:"required"`
}

// ----------
//   UPDATE
// ----------

type UpdateUserParams struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Name   string    `json:"name"`
}

type UpdateCategoryParams struct {
	UserId  uuid.UUID `json:"user_id" validate:"required"`
	CategId uuid.UUID `json:"categ_id" validate:"required"`
	Name    string    `json:"name"`
}

type UpdateFileParams struct {
	UserId   uuid.UUID `json:"user_id" validate:"required"`
	CategId  uuid.UUID `json:"categ_id" validate:"required"`
	FileId   uuid.UUID `json:"file_id" validate:"required"`
	Name     string    `json:"name"`
	FileType string    `json:"file_type"`
	Content  []byte    `json:"content"`
}

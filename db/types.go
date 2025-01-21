package db

import (
	"github.com/google/uuid"
)

type File struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Filename  string    `json:"filename"`
	UpdatedAt int64     `json:"updated_at"`
}

type Category struct {
	Id        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	DirName   string             `json:"dirname"`
	Files     map[uuid.UUID]File `json:"files"`
	UpdatedAt int64              `json:"updated_at"`
}

type User struct {
	Id         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	DirName    string                 `json:"dirname"`
	Categories map[uuid.UUID]Category `json:"categories"`
	UpdatedAt  int64                  `json:"updated_at"`
}

type Repo struct {
	Users     map[uuid.UUID]User `json:"users"`
	UpdatedAt int64              `json:"updated_at"`
}

// Estrutura para buscas

type QueryData struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt int64     `json:"updated_at"`
}

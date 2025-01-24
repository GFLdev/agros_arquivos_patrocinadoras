package handlers

import "github.com/google/uuid"

// ---------------
//   Requisições
// ---------------

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateCategoryReq struct {
	Name string `json:"name" validate:"required"`
}

type CreateFileReq struct {
	Name      string `json:"name" validate:"required"`
	Extension string `json:"extension" validate:"required"`
	Mimetype  string `json:"mimetype" validate:"required"`
	Content   []byte `json:"content" validate:"required"`
}

type UpdateUserReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateCategoryReq struct {
	UserId uuid.UUID `json:"userId" validate:"required"`
	Name   string    `json:"name" validate:"required"`
}

type UpdateFileReq struct {
	CategId   uuid.UUID `json:"categId" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Extension string    `json:"extension" validate:"required"`
	Mimetype  string    `json:"mimetype" validate:"required"`
	Content   []byte    `json:"content" validate:"required"`
}

// -------------
//   Respostas
// -------------

type ErrorRes struct {
	Message string `json:"message" validate:"required"`
	Error   string `json:"error" validate:"required"`
}

type GenericRes struct {
	Message string `json:"message" validate:"required"`
}

type LoginRes struct {
	User          string `json:"user" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
}

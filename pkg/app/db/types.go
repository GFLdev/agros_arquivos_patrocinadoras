package db

import "github.com/google/uuid"

type UserCreation struct {
	Name     string
	Password string
}

type CategCreation struct {
	UserId uuid.UUID
	Name   string
}

type FileCreation struct {
	UserId    uuid.UUID
	CategId   uuid.UUID
	Name      string
	Extension string
	Mimetype  string
	Content   []byte
}

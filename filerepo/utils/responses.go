package utils

type ErrorRes struct {
	Message string `json:"message" validate:"required"`
	Error   error  `json:"error" validate:"required"`
}

type GenericRes struct {
	Message string `json:"message" validate:"required"`
}

type LoginRes struct {
	User          string `json:"user" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
}

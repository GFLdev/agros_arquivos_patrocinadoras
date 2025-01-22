package utils

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type NameInputReq struct {
	Name string `json:"name" validate:"required"`
}

type FileInputReq struct {
	NameInputReq
	FileType  string `json:"file_type" validate:"required"`
	Extension string `json:"extension" validate:"required"`
	Content   []byte `json:"content" validate:"required"`
}

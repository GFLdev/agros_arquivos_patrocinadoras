package db

type UserCreation struct {
	Name     string
	Password string
}

//	Name      string    `json:"name" validate:"required"`
//	Extension string    `json:"extension" validate:"required"`
//	Mimetype  string    `json:"mimetype" validate:"required"`
//	Content   []byte    `json:"content" validate:"required"`
//}

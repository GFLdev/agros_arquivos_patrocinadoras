package db

type UserModel struct {
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	UpdatedAt int64  `json:"updated_at"`
}

type CategoryModel struct {
	CategId   string `json:"categ_id"`
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	UpdatedAt int64  `json:"updated_at"`
}

type FileModel struct {
	FileId    string `json:"file_id"`
	CategId   string `json:"categ_id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Mimetype  string `json:"mimetype"`
	UpdatedAt int64  `json:"updated_at"`
}

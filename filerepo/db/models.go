package db

type UserModel struct {
	UserId    string
	Name      string
	Password  string
	UpdatedAt int64
}

type CategoryModel struct {
	CategId   string
	UserId    string
	Name      string
	UpdatedAt int64
}

type FileModel struct {
	FileId    string
	CategId   string
	Name      string
	Extension string
	Mimetype  string
	UpdatedAt int64
}

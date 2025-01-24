package config

type Config struct {
	Environment string   `json:"environment" validate:"required"`
	Origins     []string `json:"origins" validate:"required"`
	Port        int      `json:"port" validate:"required"`
	Database    Database `json:"database" validate:"required"`
	JwtSecret   string   `json:"jwt_secret" validate:"required"`
	JwtExpires  int      `json:"jwt_expires" validate:"required"`
	CertFile    string   `json:"cert_file"`
	KeyFile     string   `json:"key_file"`
}

type Database struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Server   string `json:"server"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Schema   Schema `json:"schema"`
}

type Schema struct {
	Name       string            `json:"name"`
	UserTable  Table[UserTable]  `json:"user_table"`
	CategTable Table[CategTable] `json:"categ_table"`
	FileTable  Table[FileTable]  `json:"file_table"`
}

type Table[T interface{}] struct {
	Name    string `json:"name"`
	Columns T      `json:"columns"`
}

type UserTable struct {
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	UpdatedAt string `json:"updated_at"`
}

type CategTable struct {
	CategId   string `json:"categ_id"`
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

type FileTable struct {
	FileId    string `json:"file_id"`
	CategId   string `json:"categ_id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Mimetype  string `json:"mimetype"`
	UpdatedAt string `json:"updated_at"`
}

package config

import "agros_arquivos_patrocinadoras/pkg/types/db"

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
	Name       string                  `json:"name"`
	UserTable  Table[db.UserModel]     `json:"user_table"`
	CategTable Table[db.CategoryModel] `json:"categ_table"`
	FileTable  Table[db.FileModel]     `json:"file_table"`
}

type Table[T interface{}] struct {
	Name    string `json:"name"`
	Columns T      `json:"columns"`
}

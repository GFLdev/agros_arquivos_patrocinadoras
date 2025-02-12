// Package config define as estruturas de configuração utilizadas pela aplicação,
// incluindo variáveis de ambiente, configurações de origem, portas, credenciais
// de banco de dados, tabelas e outros parâmetros necessários para o funcionamento
// correto da aplicação.
package config

// Config representa a configuração principal da aplicação.
type Config struct {
	// Environment define o ambiente da aplicação (ex.: "production").
	Environment string `json:"environment" validate:"required"`
	// Origins contém a lista de origens permitidas para acessar os recursos.
	Origins []string `json:"origins" validate:"required"`
	// AdminName especifica o nome de usuário administrador no banco.
	AdminUsername string `json:"admin_username" validate:"required"`
	// AdminName especifica o nome do administrador no banco.
	AdminName string `json:"admin_name" validate:"required"`
	// Port define a porta onde o servidor da aplicação será executado.
	Port int `json:"port" validate:"required"`
	// Database armazena as configurações de conexão e esquema do banco de dados.
	Database Database `json:"database" validate:"required"`
	// JwtSecret define a chave secreta usada para geração e validação de tokens JWT.
	JwtSecret string `json:"jwt_secret" validate:"required"`
	// JwtExpires define, em minutos, o tempo de expiração para o token JWT.
	JwtExpires int `json:"jwt_expires" validate:"required"`
	// CertFile contém o caminho para o arquivo de certificado SSL (opcional).
	CertFile string `json:"cert_file"`
	// KeyFile contém o caminho para o arquivo de chave SSL (opcional).
	KeyFile string `json:"key_file"`
}

// Database representa as configurações de conexão e credenciais do banco de
// dados.
type Database struct {
	// Service define o nome do serviço do banco de dados (ex.: "ORCL").
	Service string `json:"service" validate:"required"`
	// Username define o nome do usuário usado para autenticação no banco.
	Username string `json:"username" validate:"required"`
	// Server define o endereço do servidor onde o banco de dados está hospedado.
	Server string `json:"server" validate:"required"`
	// Port define a porta de conexão para o banco de dados.
	Port string `json:"port" validate:"required"`
	// Password define a senha do usuário usada para autenticação no banco.
	Password string `json:"password" validate:"required"`
	// Schema representa as configurações e tabelas do esquema do banco de dados.
	Schema Schema `json:"schema" validate:"required"`
}

// Schema define o esquema usado no banco de dados.
type Schema struct {
	// Name define o nome do esquema no banco de dados.
	Name string `json:"name" validate:"required"`
	// UserTable representa a configuração da tabela de usuários no esquema.
	UserTable Table[UserTable] `json:"user_table" validate:"required"`
	// CategTable representa a configuração da tabela de categorias no esquema.
	CategTable Table[CategTable] `json:"categ_table" validate:"required"`
	// FileTable representa a configuração da tabela de arquivos no esquema.
	FileTable Table[FileTable] `json:"file_table" validate:"required"`
}

// Table representa uma tabela genérica usada no esquema do banco de dados.
type Table[T interface{}] struct {
	// Name define o nome da tabela no banco de dados.
	Name string `json:"name" validate:"required"`
	// Columns representa as colunas específicas associadas à tabela.
	Columns T `json:"columns" validate:"required"`
}

// UserTable representa a estrutura das colunas na tabela de usuários do banco.
type UserTable struct {
	// UserId define a coluna do identificador único de um usuário.
	UserId string `json:"user_id" validate:"required"`
	// Username define a coluna do nome de usuário, de um usuário.
	Username string `json:"username" validate:"required"`
	// Name define a coluna do nome de apresentação de um usuário.
	Name string `json:"name" validate:"required"`
	// Password define a coluna da senha de um usuário.
	Password string `json:"password" validate:"required"`
	// UpdatedAt define a coluna da última atualização do usuário.
	UpdatedAt string `json:"updated_at" validate:"required"`
}

// CategTable representa a estrutura das colunas na tabela de categorias do banco.
type CategTable struct {
	// CategId define a coluna do identificador único de uma categoria.
	CategId string `json:"categ_id" validate:"required"`
	// UserId define a coluna que referencia o identificador de um usuário.
	UserId string `json:"user_id" validate:"required"`
	// Name define a coluna do nome da categoria.
	Name string `json:"name" validate:"required"`
	// UpdatedAt define a coluna da última atualização da categoria.
	UpdatedAt string `json:"updated_at" validate:"required"`
}

// FileTable representa a estrutura das colunas na tabela de arquivos do banco.
type FileTable struct {
	// FileId define a coluna do identificador único de um arquivo.
	FileId string `json:"file_id" validate:"required"`
	// CategId define a coluna que referencia o identificador de uma categoria.
	CategId string `json:"categ_id" validate:"required"`
	// Name define a coluna do nome de um arquivo.
	Name string `json:"name" validate:"required"`
	// Extension define a coluna da extensão do arquivo (ex.: ".txt").
	Extension string `json:"extension" validate:"required"`
	// Mimetype define a coluna que especifica o tipo MIME do arquivo.
	Mimetype string `json:"mimetype" validate:"required"`
	// Blob define a coluna que especifica o conteúdo do arquivo.
	Blob string `json:"blob" validate:"required"`
	// UpdatedAt define a coluna da última atualização do arquivo.
	UpdatedAt string `json:"updated_at" validate:"required"`
}

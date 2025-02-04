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
	Service string `json:"service"`
	// Username define o nome do usuário usado para autenticação no banco.
	Username string `json:"username"`
	// Server define o endereço do servidor onde o banco de dados está hospedado.
	Server string `json:"server"`
	// Port define a porta de conexão para o banco de dados.
	Port string `json:"port"`
	// Password define a senha do usuário usada para autenticação no banco.
	Password string `json:"password"`
	// Schema representa as configurações e tabelas do esquema do banco de dados.
	Schema Schema `json:"schema"`
}

// Schema define o esquema usado no banco de dados.
type Schema struct {
	// Name define o nome do esquema no banco de dados.
	Name string `json:"name"`
	// UserTable representa a configuração da tabela de usuários no esquema.
	UserTable Table[UserTable] `json:"user_table"`
	// CategTable representa a configuração da tabela de categorias no esquema.
	CategTable Table[CategTable] `json:"categ_table"`
	// FileTable representa a configuração da tabela de arquivos no esquema.
	FileTable Table[FileTable] `json:"file_table"`
}

// Table representa uma tabela genérica usada no esquema do banco de dados.
type Table[T interface{}] struct {
	// Name define o nome da tabela no banco de dados.
	Name string `json:"name"`
	// Columns representa as colunas específicas associadas à tabela.
	Columns T `json:"columns"`
}

// UserTable representa a estrutura das colunas na tabela de usuários do banco.
type UserTable struct {
	// UserId define a coluna do identificador único de um usuário.
	UserId string `json:"user_id"`
	// Name define a coluna do nome de um usuário.
	Name string `json:"name"`
	// Password define a coluna da senha de um usuário.
	Password string `json:"password"`
	// UpdatedAt define a coluna da última atualização do usuário.
	UpdatedAt string `json:"updated_at"`
}

// CategTable representa a estrutura das colunas na tabela de categorias do banco.
type CategTable struct {
	// CategId define a coluna do identificador único de uma categoria.
	CategId string `json:"categ_id"`
	// UserId define a coluna que referencia o identificador de um usuário.
	UserId string `json:"user_id"`
	// Name define a coluna do nome da categoria.
	Name string `json:"name"`
	// UpdatedAt define a coluna da última atualização da categoria.
	UpdatedAt string `json:"updated_at"`
}

// FileTable representa a estrutura das colunas na tabela de arquivos do banco.
type FileTable struct {
	// FileId define a coluna do identificador único de um arquivo.
	FileId string `json:"file_id"`
	// CategId define a coluna que referencia o identificador de uma categoria.
	CategId string `json:"categ_id"`
	// Name define a coluna do nome de um arquivo.
	Name string `json:"name"`
	// Extension define a coluna da extensão do arquivo (ex.: ".txt").
	Extension string `json:"extension"`
	// Mimetype define a coluna que especifica o tipo MIME do arquivo.
	Mimetype string `json:"mimetype"`
	// Blob define a coluna que especifica o conteúdo do arquivo.
	Blob string `json:"blob"`
	// UpdatedAt define a coluna da última atualização do arquivo.
	UpdatedAt string `json:"updated_at"`
}

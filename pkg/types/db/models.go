// Package db define os modelos de dados para a aplicação. Ele inclui estruturas
// para representar usuários, categorias e arquivos, com atributos relevantes
// para identificação e rastreamento de atualizações de dados.
package db

// UserModel representa o modelo do usuário armazenado no banco de dados.
type UserModel struct {
	// UserId é um identificador único para o usuário.
	UserId string `json:"user_id"`
	// Username representa o nome de usuário, do usuário.
	Username string `json:"username"`
	// Name representa o nome de apresentação do usuário.
	Name string `json:"name"`
	// Password armazena a senha (hash) do usuário.
	Password string `json:"password"`
	// UpdatedAt representa o timestamp da última atualização dos dados
	// do usuário, armazenado como um tempo Unix em segundos.
	UpdatedAt int64 `json:"updated_at"`
}

// CategModel representa o modelo da categoria armazenada no banco de dados.
type CategModel struct {
	// CategId é o identificador único de uma categoria.
	CategId string `json:"categ_id"`
	// UserId é o identificador único de um usuário associado à categoria.
	UserId string `json:"user_id"`
	// Name é o nome da categoria.
	Name string `json:"name"`
	// UpdatedAt representa o timestamp da última atualização dos dados
	// da categoria, armazenado como um tempo Unix em segundos.
	UpdatedAt int64 `json:"updated_at"`
}

// FileModel representa o modelo do arquivo armazenado no banco de dados.
type FileModel struct {
	// FileId representa o identificador único do arquivo.
	FileId string `json:"file_id"`
	// CategId representa o identificador único da categoria associada ao
	// arquivo.
	CategId string `json:"categ_id"`
	// Name representa o nome do arquivo.
	Name string `json:"name"`
	// Extension especifica a extensão do arquivo (por exemplo, ".txt", ".jpg").
	Extension string `json:"extension"`
	// Mimetype representa o tipo de mídia do arquivo, indicando seu formato e
	// como ele deve ser interpretado ou processado.
	Mimetype string `json:"mimetype"`
	// Blob armazena os dados brutos do arquivo como uma sequência de bytes.
	Blob []byte `json:"blob"`
	// UpdatedAt representa o timestamp da última atualização dos dados
	// do arquivo, armazenado como um tempo Unix em segundos.
	UpdatedAt int64 `json:"updated_at"`
}

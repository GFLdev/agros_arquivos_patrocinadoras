package fs

import (
	"agros_arquivos_patrocinadoras/filerepo/services/logger"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

// CreateUser cria um novo usuário no repositório.
func (fs *FS) CreateUser(p CreateUserParams) error {
	// Cria nova instância de usuário
	id := uuid.New()
	ts := time.Now().Unix()
	newUser := User{
		Id:         id,
		Name:       p.Name,
		Categories: make(map[uuid.UUID]Category),
		UpdatedAt:  ts,
	}
	path := fmt.Sprintf("fs/user_%s", id.String())

	// Cria o diretório deste usuário
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar pasta do usuário: %v", err)
	}
	newUser.Path = path

	// Adiciona usuário ao repositório
	fs.Users[id] = newUser
	fs.UpdatedAt = ts

	// Salvar em disco
	return StructToFile[FS]("fs/track.json",
		fs,
		logger.CreateLogger(),
	)
}

// CreateCategory cria uma nova categoria associada a um usuário.
func (fs *FS) CreateCategory(p CreateCategoryParams) error {
	// Verifica de existência do usuário
	user, ok := fs.Users[p.UserId]
	if !ok {
		return fmt.Errorf("usuário não encontrado no repositório")
	}

	// Cria nova instância de categoria
	id := uuid.New()
	ts := time.Now().Unix()
	newCategory := Category{
		Id:        id,
		Name:      p.Name,
		Files:     make(map[uuid.UUID]File),
		UpdatedAt: ts,
	}
	path := fmt.Sprintf("%s/categ_%s", user.Path, id.String())

	// Cria o diretório da nova categoria
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar pasta da categoria: %v", err)
	}
	newCategory.Path = path

	// Adiciona categoria ao repositório
	user.Categories[id] = newCategory
	user.UpdatedAt = ts
	fs.Users[p.UserId] = user
	fs.UpdatedAt = ts

	// Salvar em disco
	return StructToFile[FS]("fs/track.json",
		fs,
		logger.CreateLogger(),
	)
}

// CreateFile cria um novo arquivo numa categoria associada a um usuário.
func (fs *FS) CreateFile(p CreateFileParams) error {
	// Verifica de existência do usuário
	user, ok := fs.Users[p.UserId]
	if !ok {
		return fmt.Errorf("usuário não encontrado no repositório")
	}

	// Verifica de existência da categoria
	categ, ok := user.Categories[p.CategId]
	if !ok {
		return fmt.Errorf("categoria não encontrada no repositório")
	}

	// Cria nova instância de arquivo
	id := uuid.New()
	ts := time.Now().Unix()
	newFile := File{
		Id:        id,
		Name:      p.Name,
		FileType:  p.FileType,
		Extension: p.Extension,
		UpdatedAt: ts,
	}
	path := fmt.Sprintf("%s/file_%s", categ.Path, id)

	// Salva o conteúdo deste arquivo em disco
	err := WriteToFile(path, p.Content, logger.CreateLogger())
	if err != nil {
		return err
	}
	newFile.Path = path

	// Adiciona arquivo ao repositório
	// Se for o primeiro
	if len(categ.Files) == 0 {
		categ.Files = make(map[uuid.UUID]File)
	}
	categ.Files[id] = newFile
	categ.UpdatedAt = ts
	user.Categories[p.CategId] = categ
	user.UpdatedAt = ts
	fs.Users[p.UserId] = user
	fs.UpdatedAt = ts

	// Salvar em disco
	return StructToFile[FS]("fs/track.json",
		fs,
		logger.CreateLogger(),
	)
}

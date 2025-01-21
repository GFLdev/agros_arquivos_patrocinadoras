package db

import (
	"agros_arquivos_patrocinadoras/clogger"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

func newQuery(id uuid.UUID, name string, ts int64) QueryData {
	return QueryData{
		Id:        id,
		Name:      name,
		UpdatedAt: ts,
	}
}

// CREATE

func (repo *Repo) CreateUser(name string) error {
	// Cria nova instância de usuário
	id := uuid.New()
	ts := time.Now().Unix()
	newUser := User{
		Id:         id,
		Name:       name,
		Categories: make(map[uuid.UUID]Category),
		UpdatedAt:  ts,
	}
	dirName := fmt.Sprintf("repo/user_%s", id.String())

	// Cria o diretório deste usuário
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar pasta do usuário: %v", err)
	}
	newUser.DirName = dirName

	// Adiciona usuário ao repositório
	repo.Users[id] = newUser
	repo.UpdatedAt = ts
	return nil
}

func (repo *Repo) CreateCategory(userId uuid.UUID, name string) error {
	// Verifica de existência do usuário
	user, ok := repo.Users[userId]
	if !ok {
		return fmt.Errorf("usuário não encontrado no repositório")
	}

	// Cria nova instância de categoria
	id := uuid.New()
	ts := time.Now().Unix()
	newCategory := Category{
		Id:        id,
		Name:      name,
		Files:     make(map[uuid.UUID]File),
		UpdatedAt: ts,
	}
	dirName := fmt.Sprintf("%s/categ_%s", user.DirName, id.String())

	// Cria o diretório da nova categoria
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("erro ao criar pasta da categoria: %v", err)
	}
	newCategory.DirName = dirName

	// Adiciona categoria ao repositório
	user.Categories[id] = newCategory
	user.UpdatedAt = ts
	repo.Users[userId] = user
	repo.UpdatedAt = ts
	return nil
}

func (repo *Repo) CreateFile(
	userId uuid.UUID,
	categId uuid.UUID,
	name string,
	content []byte,
) error {
	// Verifica de existência do usuário
	user, ok := repo.Users[userId]
	if !ok {
		return fmt.Errorf("usuário não encontrado no repositório")
	}

	// Verifica de existência da categoria
	categ, ok := user.Categories[categId]
	if !ok {
		return fmt.Errorf("categoria não encontrada no repositório")
	}

	// Cria nova instância de arquivo
	id := uuid.New()
	ts := time.Now().Unix()
	newFile := File{
		Id:        id,
		Name:      name,
		UpdatedAt: ts,
	}
	basename := fmt.Sprintf("%s/file_%s", categ.DirName, id)

	// Salva o conteúdo deste arquivo em disco
	err := WriteToFile(basename, content, clogger.CreateLogger())
	if err != nil {
		return err
	}
	newFile.Basename = basename

	// Adiciona arquivo ao repositório
	// Se for o primeiro
	if len(categ.Files) == 0 {
		categ.Files = make(map[uuid.UUID]File)
	}
	categ.Files[id] = newFile
	categ.UpdatedAt = ts
	user.Categories[categId] = categ
	user.UpdatedAt = ts
	repo.Users[userId] = user
	repo.UpdatedAt = ts
	return nil
}

// READ

func (repo *Repo) GetUserById(userId uuid.UUID) QueryData {
	user, ok := repo.Users[userId]
	if !ok {
		return QueryData{}
	}
	return QueryData{user.Id, user.Name, user.UpdatedAt}
}

func (repo *Repo) GetAllUsers() []QueryData {
	users := make([]QueryData, len(repo.Users))
	for _, user := range repo.Users {
		users = append(users, QueryData{user.Id, user.Name, user.UpdatedAt})
	}
	return users
}

func (repo *Repo) GetCategoryById(
	userId uuid.UUID,
	categId uuid.UUID,
) QueryData {
	user, ok := repo.Users[userId]
	if !ok {
		return QueryData{}
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return QueryData{}
	}
	return QueryData{categ.Id, categ.Name, categ.UpdatedAt}
}

func (repo *Repo) GetAllCategories(userId uuid.UUID) []QueryData {
	user, ok := repo.Users[userId]
	if !ok {
		return []QueryData{}
	}

	categs := make([]QueryData, len(user.Categories))
	for _, categ := range user.Categories {
		categs = append(categs, QueryData{categ.Id, categ.Name, categ.UpdatedAt})
	}
	return categs
}

func (repo *Repo) GetFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) QueryData {
	user, ok := repo.Users[userId]
	if !ok {
		return QueryData{}
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return QueryData{}
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return QueryData{}
	}
	return QueryData{file.Id, file.Name, file.UpdatedAt}
}

func (repo *Repo) GetAllFiles(userId uuid.UUID, categId uuid.UUID) []QueryData {
	user, ok := repo.Users[userId]
	if !ok {
		return []QueryData{}
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return []QueryData{}
	}

	files := make([]QueryData, len(categ.Files))
	for _, file := range categ.Files {
		files = append(files, QueryData{file.Id, file.Name, file.UpdatedAt})
	}
	return files
}

// UPDATE

func (repo *Repo) UpdateUserById(userId uuid.UUID, name string) bool {
	user, ok := repo.Users[userId]
	if !ok {
		return false
	}

	ts := time.Now().Unix()
	user.Name = name
	user.UpdatedAt = ts
	repo.Users[userId] = user
	repo.UpdatedAt = ts

	return true
}

func (repo *Repo) UpdateCategoryById(
	userId uuid.UUID,
	categId uuid.UUID,
	name string,
) bool {
	user, ok := repo.Users[userId]
	if !ok {
		return false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return false
	}

	ts := time.Now().Unix()
	categ.Name = name
	categ.UpdatedAt = ts
	user.Categories[categId] = categ
	user.UpdatedAt = ts
	repo.Users[userId] = user
	repo.UpdatedAt = ts

	return true
}

func (repo *Repo) UpdateFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
	name string,
) bool {
	user, ok := repo.Users[userId]
	if !ok {
		return false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return false
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return false
	}

	ts := time.Now().Unix()
	file.Name = name
	file.UpdatedAt = ts
	categ.Files[fileId] = file
	categ.UpdatedAt = ts
	user.Categories[categId] = categ
	user.UpdatedAt = ts
	repo.Users[userId] = user
	repo.UpdatedAt = ts
	return true
}

// DELETE

func (repo *Repo) DeleteUserById(userId uuid.UUID) bool {
	_, ok := repo.Users[userId]
	if ok {
		delete(repo.Users, userId)
		return true
	}
	return false
}

func (repo *Repo) DeleteCategoryById(
	userId uuid.UUID,
	categId uuid.UUID,
) bool {
	user, ok := repo.Users[userId]
	if !ok {
		return false
	}

	_, ok = user.Categories[categId]
	if !ok {
		return false
	}

	delete(repo.Users[userId].Categories, categId)
	return true
}

func (repo *Repo) DeleteFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) bool {
	user, ok := repo.Users[userId]
	if !ok {
		return false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return false
	}

	_, ok = categ.Files[fileId]
	if !ok {
		return false
	}

	delete(repo.Users[userId].Categories[categId].Files, fileId)
	return true
}

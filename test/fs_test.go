package test

import (
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"github.com/google/uuid"
	"path/filepath"
	"testing"
)

type mockData struct {
	Fs      fs.FileSystem
	UserId  uuid.UUID
	CategId uuid.UUID
	FileId  uuid.UUID
	Content []byte
}

func newMock() *mockData {
	return &mockData{
		Fs:      fs.FileSystem{Root: TestRoot},
		UserId:  uuid.New(),
		CategId: uuid.New(),
		FileId:  uuid.New(),
		Content: []byte("Hello World"),
	}
}

func TestFileSystem_CreateEntity(t *testing.T) {
	// Criação
	mock := newMock()
	userPath := filepath.Join(mock.Fs.Root, mock.UserId.String())
	categPath := filepath.Join(userPath, mock.CategId.String())
	filePath := filepath.Join(categPath, mock.FileId.String()+".txt")

	// Testes
	if err := mock.Fs.CreateEntity(userPath, nil, fs.User); err != nil {
		t.Error(err)
	}
	if err := mock.Fs.CreateEntity(categPath, nil, fs.Category); err != nil {
		t.Error(err)
	}
	if err := mock.Fs.CreateEntity(filePath, &mock.Content, fs.File); err != nil {
		t.Error(err)
	}
}

func TestFileSystem_GetEntity(t *testing.T) {
	// Criação
	mock := newMock()
	userPath := filepath.Join(mock.Fs.Root, mock.UserId.String())
	categPath := filepath.Join(userPath, mock.CategId.String())
	filePath := filepath.Join(categPath, mock.FileId.String()+".txt")
	if err := mock.Fs.CreateEntity(userPath, nil, fs.User); err != nil {
		panic(err)
	}
	if err := mock.Fs.CreateEntity(categPath, nil, fs.Category); err != nil {
		panic(err)
	}
	if err := mock.Fs.CreateEntity(filePath, &mock.Content, fs.File); err != nil {
		panic(err)
	}

	// Testes
	if !mock.Fs.EntityExists(userPath) {
		t.Error("Usuário não existe")
	}
	if !mock.Fs.EntityExists(categPath) {
		t.Error("Categoria não existe")
	}
	if !mock.Fs.EntityExists(filePath) {
		t.Error("Arquivo não existe")
	}
}

func TestFileSystem_UpdateEntity(t *testing.T) {
	// Criação
	mock := newMock()
	userPath := filepath.Join(mock.Fs.Root, mock.UserId.String())
	categPath := filepath.Join(userPath, mock.CategId.String())
	filePath := filepath.Join(categPath, mock.FileId.String()+".txt")
	if err := mock.Fs.CreateEntity(userPath, nil, fs.User); err != nil {
		panic(err)
	}
	if err := mock.Fs.CreateEntity(categPath, nil, fs.Category); err != nil {
		panic(err)
	}
	if err := mock.Fs.CreateEntity(filePath, &mock.Content, fs.File); err != nil {
		panic(err)
	}

	// Testes
	newCateg := uuid.New()
	newCategPath := filepath.Join(userPath, newCateg.String())
	if err := mock.Fs.CreateEntity(newCategPath, nil, fs.Category); err != nil {
		panic(err)
	}
	fileNewPath := filepath.Join(newCategPath, mock.FileId.String()+".txt")
	if err := mock.Fs.UpdateEntity(filePath, fileNewPath); err != nil {
		t.Error(err)
	}

	newUserId := uuid.New()
	newUserPath := filepath.Join(mock.Fs.Root, newUserId.String())
	if err := mock.Fs.CreateEntity(newUserPath, nil, fs.User); err != nil {
		panic(err)
	}
	categNewPath := filepath.Join(newUserPath, mock.CategId.String())
	if err := mock.Fs.UpdateEntity(categPath, categNewPath); err != nil {
		t.Error(err)
	}
}

func TestFileSystem_DeleteEntity(t *testing.T) {
	// Criação
	mock := newMock()
	userPath := filepath.Join(mock.Fs.Root, mock.UserId.String())
	categPath := filepath.Join(userPath, mock.CategId.String())
	filePath := filepath.Join(categPath, mock.FileId.String()+".txt")
	if err := mock.Fs.CreateEntity(userPath, nil, fs.User); err != nil {
		panic(err)
	}
	if err := mock.Fs.CreateEntity(categPath, nil, fs.Category); err != nil {
		panic(err)
	}
	if err := mock.Fs.CreateEntity(filePath, &mock.Content, fs.File); err != nil {
		panic(err)
	}

	// Testes
	// Quebra de parente-filho
	if err := mock.Fs.DeleteEntity(userPath); err == nil {
		t.Error("usuário excluído sem excluir as categorias")
	}
	if err := mock.Fs.DeleteEntity(categPath); err == nil {
		t.Error("categoria excluída sem excluir os arquivos")
	}
	// Corretos
	if err := mock.Fs.DeleteEntity(filePath); err != nil {
		t.Error(err)
	}
	if err := mock.Fs.DeleteEntity(categPath); err != nil {
		t.Error(err)
	}
	if err := mock.Fs.DeleteEntity(userPath); err != nil {
		t.Error(err)
	}
}

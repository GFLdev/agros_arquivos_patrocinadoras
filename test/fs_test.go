package test

import (
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
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

func TestFileSystem_CreateReadEntities(t *testing.T) {
	// Cenário positivo
	t.Run(
		"Should_Create_And_Read_Entities",
		func(t *testing.T) {
			m := newMock()

			// Caminhos
			userPath := filepath.Join(m.Fs.Root, m.UserId.String())
			categPath := filepath.Join(userPath, m.CategId.String())
			filePath := filepath.Join(categPath, m.FileId.String()+".txt")
			content := []byte("Hello World")

			// Criar usuário
			err := m.Fs.CreateEntity(userPath, nil, fs.User)
			assert.NoError(t, err)
			ok := m.Fs.EntityExists(userPath)
			assert.True(t, ok, "Usuário não foi criado no caminho "+userPath)

			// Criar categoria
			err = m.Fs.CreateEntity(categPath, nil, fs.Category)
			assert.NoError(t, err)
			ok = m.Fs.EntityExists(categPath)
			assert.True(t, ok, "Categoria não foi criada no caminho "+categPath)

			// Criar arquivo
			err = m.Fs.CreateEntity(filePath, &content, fs.File)
			assert.NoError(t, err)
			ok = m.Fs.EntityExists(filePath)
			assert.True(t, ok, "Usuário não foi criado no caminho "+filePath)

			// Verificar conteúdo do arquivo
			file, err := os.Open(filePath)
			assert.NoError(t, err)
			writtenContent, err := io.ReadAll(file)
			assert.NoError(t, err)
			assert.Equal(t, content, writtenContent)
		},
	)

	// Cenários negativos
	t.Run(
		"Should_Return_Error_When_Creating_Entity_With_Invalid_Path",
		func(t *testing.T) {
			m := newMock()
			invalidPath := "caminho_invalido"

			// Tentativa de criar entidade num caminho inválido
			err := m.Fs.CreateEntity(invalidPath, nil, fs.User)
			assert.Error(t, err)
		},
	)

	t.Run(
		"Should_Return_Error_When_Creating_Entity_With_Non_Existing_Parent",
		func(t *testing.T) {
			m := newMock()
			userPath := filepath.Join(m.Fs.Root, m.UserId.String())
			categPath := filepath.Join(userPath, m.CategId.String())

			// Tentativa de criar entidade sem ter pai
			err := m.Fs.CreateEntity(categPath, nil, fs.Category)
			assert.Error(t, err)
		},
	)
}

func TestFileSystem_UpdateEntity(t *testing.T) {
	// Configuração inicial do mock e caminhos
	mock := newMock()
	userPath := filepath.Join(mock.Fs.Root, mock.UserId.String())
	categPath := filepath.Join(userPath, mock.CategId.String())
	filePath := filepath.Join(categPath, mock.FileId.String()+".txt")

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

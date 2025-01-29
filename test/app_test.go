package test

import (
	"agros_arquivos_patrocinadoras/pkg/app"
	"agros_arquivos_patrocinadoras/pkg/app/config"
	"agros_arquivos_patrocinadoras/pkg/app/context"
	"agros_arquivos_patrocinadoras/pkg/app/db"
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/app/logger"
	"github.com/google/uuid"
	"testing"
)

func newMockUser() app.UserParams {
	return app.UserParams{
		Name:     "test",
		Password: "test123456789",
	}
}

func newMockCateg(userId uuid.UUID) app.CategParams {
	return app.CategParams{
		UserId: userId,
		Name:   "test",
	}
}

func newMockFile(userId, categId uuid.UUID) app.FileParams {
	content := []byte("Hello World")
	return app.FileParams{
		UserId:    userId,
		CategId:   categId,
		Name:      "test",
		Extension: ".txt",
		Mimetype:  "text/plain",
		Content:   &content,
	}
}

func newContext() *context.Context {
	logr := logger.CreateLogger()
	cfg := config.LoadConfig(logr)
	filesystem := &fs.FileSystem{Root: TestRoot}
	dataBase := db.GetSqlDB(&cfg.Database, logr)
	return &context.Context{
		Logger:     logr,
		Config:     cfg,
		FileSystem: filesystem,
		DB:         dataBase,
	}
}

func TestApp_CreateUser(t *testing.T) {
	ctx := newContext()

	// Criar usuário
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		t.Error(err)
	}
	if userId == uuid.Nil {
		t.Error("Id do usuário nulo")
	}

	// Obter usuário e validar dados obtidos
	user, err := app.QueryUserById(ctx, userId)
	if err != nil {
		t.Error(err)
	}
	if user.Name != userParams.Name {
		t.Error("Nomes do usuário diferentes")
	}

	// Autenticar usuário
	credentials, err := app.GetCredentials(ctx, userParams)
	if err != nil {
		t.Error(err)
	}
	if credentials != userId {
		t.Error("Ids do usuário diferente")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		t.Error(err)
	}
}

func TestApp_CreateCategory(t *testing.T) {
	ctx := newContext()

	// Criar usuário e categoria
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		panic(err)
	} else if userId == uuid.Nil {
		panic("usuário com Id nulo")
	}
	categParams := newMockCateg(userId)
	categId, err := app.CreateCategory(ctx, categParams)
	if err != nil {
		t.Error(err)
	}
	if categId == uuid.Nil {
		t.Error("Id da categoria nulo")
	}

	// Obter categoria e validar dados obtidos
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil {
		t.Error(err)
	}
	if categ.UserId != categParams.UserId.String() {
		t.Error("Id de usuário, da categoria, diferente")
	}
	if categ.Name != categParams.Name {
		t.Error("Nomes da categoria diferentes")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, delUser); err == nil {
		t.Error("Deleção de usuário com entidades filhas")
	}
	delCateg := app.CategDelete{UserId: userId, CategId: categId}
	if err = app.DeleteCategory(ctx, delCateg); err != nil {
		t.Error(err)
	}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		panic(err)
	}
}

func TestApp_CreateFile(t *testing.T) {
	ctx := newContext()

	// Criar usuário, categoria e arquivo
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		panic(err)
	} else if userId == uuid.Nil {
		panic("usuário com Id nulo")
	}
	categParams := newMockCateg(userId)
	categId, err := app.CreateCategory(ctx, categParams)
	if err != nil {
		t.Error(err)
	} else if categId == uuid.Nil {
		panic("categoria com Id nulo")
	}
	fileParams := newMockFile(userId, categId)
	fileId, err := app.CreateFile(ctx, fileParams)
	if err != nil {
		t.Error(err)
	}
	if fileId == uuid.Nil {
		t.Error("Id do arquivo nulo")
	}

	// Obter arquivo e validar dados obtidos
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		t.Error(err)
	}
	if file.Mimetype != fileParams.Mimetype {
		t.Error("Mimetypes do arquivo diferentes")
	}
	if file.Extension != fileParams.Extension {
		t.Error("Extensões do arquivo diferentes")
	}
	if file.CategId != fileParams.CategId.String() {
		t.Error("Ids da categoria, do arquivo, diferentes")
	}
	if file.Name != fileParams.Name {
		t.Error("Nomes do arquivo diferentes")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, delUser); err == nil {
		t.Error("Deleção de usuário com entidades filhas")
	}
	delCateg := app.CategDelete{UserId: userId, CategId: categId}
	if err = app.DeleteCategory(ctx, delCateg); err == nil {
		t.Error("Deleção de categoria com entidades filhas")
	}
	delFile := app.FileDelete{
		CategDelete: delCateg,
		FileId:      fileId,
		Extension:   fileParams.Extension,
	}
	if err = app.DeleteFile(ctx, delFile); err != nil {
		t.Error(err)
	}
	if err = app.DeleteCategory(ctx, delCateg); err != nil {
		panic(err)
	}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		panic(err)
	}
}

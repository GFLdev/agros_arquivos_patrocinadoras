package test

import (
	"agros_arquivos_patrocinadoras/pkg/app"
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

func TestApp_CreateDeleteUser(t *testing.T) {
	ctx := newContext()

	// Criar usuário
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		t.Error(err)
	}
	if userId == uuid.Nil {
		t.Error("id do usuário nulo")
	}

	// Obter usuário e validar dados obtidos
	user, err := app.QueryUserById(ctx, userId)
	if err != nil {
		t.Error(err)
	}
	if user.Name != userParams.Name {
		t.Error("nomes do usuário diferentes")
	}

	// Autenticar usuário
	credentials, err := app.GetCredentials(ctx, userParams)
	if err != nil {
		t.Error(err)
	}
	if credentials != userId {
		t.Error("ids do usuário diferente")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		t.Error(err)
	}
}

func TestApp_CreateDeleteCategory(t *testing.T) {
	ctx := newContext()

	// Criar usuário e categoria
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		panic(err)
	} else if userId == uuid.Nil {
		panic("id do usuário nulo")
	}
	categParams := newMockCateg(userId)
	categId, err := app.CreateCategory(ctx, categParams)
	if err != nil {
		t.Error(err)
	}
	if categId == uuid.Nil {
		t.Error("id da categoria nulo")
	}

	// Obter categoria e validar dados obtidos
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil {
		t.Error(err)
	}
	if categ.UserId != categParams.UserId.String() {
		t.Error("id de usuário, da categoria, diferente")
	}
	if categ.Name != categParams.Name {
		t.Error("nomes da categoria diferentes")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, delUser); err == nil {
		t.Error("deleção de usuário com entidades filhas")
	}
	delCateg := app.CategDelete{UserDelete: delUser, CategId: categId}
	if err = app.DeleteCategory(ctx, delCateg); err != nil {
		t.Error(err)
	}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		panic(err)
	}
}

func TestApp_CreateDeleteFile(t *testing.T) {
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
		t.Error("id do arquivo nulo")
	}

	// Obter arquivo e validar dados obtidos
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		t.Error(err)
	}
	if file.Mimetype != fileParams.Mimetype {
		t.Error("mimetypes do arquivo diferentes")
	}
	if file.Extension != fileParams.Extension {
		t.Error("extensões do arquivo diferentes")
	}
	if file.CategId != fileParams.CategId.String() {
		t.Error("ids da categoria, do arquivo, diferentes")
	}
	if file.Name != fileParams.Name {
		t.Error("nomes do arquivo diferentes")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	delCateg := app.CategDelete{UserDelete: delUser, CategId: categId}
	delFile := app.FileDelete{
		CategDelete: delCateg,
		FileId:      fileId,
		Extension:   fileParams.Extension,
	}
	if err = app.DeleteUser(ctx, delUser); err == nil {
		t.Error("deleção de usuário com entidades filhas")
	}
	if err = app.DeleteCategory(ctx, delCateg); err == nil {
		t.Error("deleção de categoria com entidades filhas")
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

func TestApp_UpdateUser(t *testing.T) {
	ctx := newContext()

	// Criar usuário
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		panic(err)
	}
	if userId == uuid.Nil {
		panic("id do usuário nulo")
	}

	// Atualizar dados
	up := app.UserUpdate{
		UserId: userId,
		UserParams: app.UserParams{
			Name:     "test2",
			Password: "123456789test",
		},
	}
	if err = app.UpdateUser(ctx, up); err != nil {
		t.Error(err)
	}

	// Obter usuário e validar dados obtidos
	user, err := app.QueryUserById(ctx, userId)
	if err != nil {
		t.Error(err)
	}
	if user.Name != up.UserParams.Name {
		t.Error("nomes do usuário diferentes")
	}

	// Autenticar usuário
	credentials, err := app.GetCredentials(ctx, up.UserParams)
	if err != nil {
		t.Error(err)
	}
	if credentials != userId {
		t.Error("ids do usuário diferente")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		panic(err)
	}
}

func TestApp_UpdateCategory(t *testing.T) {
	ctx := newContext()

	// Criar usuários
	user0Params := newMockUser()
	user1Params := newMockUser()
	user0Id, err := app.CreateUser(ctx, user0Params)
	if err != nil {
		panic(err)
	} else if user0Id == uuid.Nil {
		panic("id do usuário nulo")
	}
	user1Id, err := app.CreateUser(ctx, user1Params)
	if err == nil {
		t.Error("nome de usuário repetido")
	}
	user1Params.Name = "test2"
	user1Id, err = app.CreateUser(ctx, user1Params)
	if err != nil {
		panic(err)
	} else if user1Id == uuid.Nil {
		panic("id do usuário nulo")
	}

	// Criar categoria
	categParams := newMockCateg(user0Id)
	categId, err := app.CreateCategory(ctx, categParams)
	if err != nil {
		panic(err)
	}
	if categId == uuid.Nil {
		panic("id da categoria nulo")
	}

	// Atualizar dados
	up := app.CategUpdate{
		CategId:   categId,
		OldUserId: categParams.UserId,
		OldName:   categParams.Name,
		CategParams: app.CategParams{
			UserId: uuid.New(),
			Name:   "test2",
		},
	}
	if err = app.UpdateCategory(ctx, up); err == nil {
		t.Error("categoria atualizada com usuário não existente")
	}
	up.CategParams.UserId = user1Id
	if err = app.UpdateCategory(ctx, up); err != nil {
		t.Error(err)
	}

	// Obter categoria e validar dados obtidos
	categ, err := app.QueryCategoryById(ctx, categId)
	if err != nil {
		t.Error(err)
	}
	if categ.UserId != up.CategParams.UserId.String() {
		t.Error("id de usuário, da categoria, diferente")
	}
	if categ.Name != up.CategParams.Name {
		t.Error("nomes da categoria diferentes")
	}

	// Deletar do banco
	del0User := app.UserDelete{UserId: user0Id}
	del1User := app.UserDelete{UserId: user1Id}
	delCateg := app.CategDelete{UserDelete: del1User, CategId: categId}
	if err = app.DeleteCategory(ctx, delCateg); err != nil {
		panic(err)
	}
	if err = app.DeleteUser(ctx, del0User); err != nil {
		panic(err)
	}
	if err = app.DeleteUser(ctx, del1User); err != nil {
		panic(err)
	}
}

func TestApp_UpdateFile(t *testing.T) {
	ctx := newContext()

	// Criar usuário
	userParams := newMockUser()
	userId, err := app.CreateUser(ctx, userParams)
	if err != nil {
		panic(err)
	} else if userId == uuid.Nil {
		panic("id do usuário nulo")
	}

	// Criar categorias
	categ0Params := newMockCateg(userId)
	categ1Params := newMockCateg(userId)
	categ0Id, err := app.CreateCategory(ctx, categ0Params)
	if err != nil {
		panic(err)
	}
	if categ0Id == uuid.Nil {
		panic("id da categoria nulo")
	}
	categ1Id, err := app.CreateCategory(ctx, categ1Params)
	if err != nil {
		panic(err)
	}
	if categ0Id == uuid.Nil {
		panic("id da categoria nulo")
	}

	// Criar arquivo
	fileParams := newMockFile(userId, categ0Id)
	fileId, err := app.CreateFile(ctx, fileParams)
	if err != nil {
		panic(err)
	}
	if fileId == uuid.Nil {
		panic("id do arquivo nulo")
	}

	// Atualizar dados
	newContent := []byte("pais,estado,cidade")
	up := app.FileUpdate{
		FileId:       fileId,
		OldCategId:   fileParams.CategId,
		OldName:      fileParams.Name,
		OldExtension: fileParams.Extension,
		OldMimetype:  fileParams.Mimetype,
		FileParams: app.FileParams{
			UserId:    uuid.New(),
			CategId:   uuid.New(),
			Name:      "test2",
			Extension: ".csv",
			Mimetype:  "text/csv",
			Content:   &newContent,
		},
	}
	if err = app.UpdateFile(ctx, up); err == nil {
		t.Error("arquivo atualizado com usuário não existente")
	}
	up.FileParams.UserId = userId
	if err = app.UpdateFile(ctx, up); err == nil {
		t.Error("arquivo atualizado com categoria não existente")
	}
	up.FileParams.CategId = categ1Id
	if err = app.UpdateFile(ctx, up); err != nil {
		t.Error(err)
	}

	// Obter arquivo e validar dados obtidos
	file, err := app.QueryFileById(ctx, fileId)
	if err != nil {
		t.Error(err)
	}
	if file.Mimetype != up.FileParams.Mimetype {
		t.Error("mimetypes do arquivo diferentes")
	}
	if file.Extension != up.FileParams.Extension {
		t.Error("extensões do arquivo diferentes")
	}
	if file.CategId != up.FileParams.CategId.String() {
		t.Error("ids da categoria, do arquivo, diferentes")
	}
	if file.Name != up.FileParams.Name {
		t.Error("nomes do arquivo diferentes")
	}

	// Deletar do banco
	delUser := app.UserDelete{UserId: userId}
	del0Categ := app.CategDelete{UserDelete: delUser, CategId: categ0Id}
	del1Categ := app.CategDelete{UserDelete: delUser, CategId: categ1Id}
	delFile := app.FileDelete{
		CategDelete: del1Categ,
		FileId:      fileId,
		Extension:   file.Extension,
	}
	if err = app.DeleteFile(ctx, delFile); err != nil {
		t.Error(err)
	}
	if err = app.DeleteCategory(ctx, del0Categ); err != nil {
		panic(err)
	}
	if err = app.DeleteCategory(ctx, del1Categ); err != nil {
		panic(err)
	}
	if err = app.DeleteUser(ctx, delUser); err != nil {
		panic(err)
	}
}

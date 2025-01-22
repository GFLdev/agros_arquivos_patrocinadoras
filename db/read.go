package db

import "github.com/google/uuid"

func (repo *Repo) GetUserById(userId uuid.UUID) (EntityData, bool) {
	user, ok := repo.Users[userId]
	if !ok {
		return EntityData{}, false
	}

	return EntityData{
		Id:        user.Id,
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}, true
}

func (repo *Repo) GetAllUsers() ([]EntityData, bool) {
	if repo.Users == nil {
		return []EntityData{}, false
	}

	var users []EntityData
	for _, user := range repo.Users {
		users = append(users, EntityData{
			Id:        user.Id,
			Name:      user.Name,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return users, true
}

func (repo *Repo) GetCategoryById(
	userId uuid.UUID,
	categId uuid.UUID,
) (EntityData, bool) {
	user, ok := repo.Users[userId]
	if !ok {
		return EntityData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return EntityData{}, false
	}

	return EntityData{
		Id:        categ.Id,
		Name:      categ.Name,
		UpdatedAt: categ.UpdatedAt,
	}, true
}

func (repo *Repo) GetAllCategories(userId uuid.UUID) ([]EntityData, bool) {
	user, ok := repo.Users[userId]
	if !ok {
		return []EntityData{}, false
	}

	var categs []EntityData
	for _, categ := range user.Categories {
		categs = append(categs, EntityData{
			Id:        categ.Id,
			Name:      categ.Name,
			UpdatedAt: categ.UpdatedAt,
		})
	}

	return categs, true
}

func (repo *Repo) GetFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) (FileData, bool) {
	user, ok := repo.Users[userId]
	if !ok {
		return FileData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return FileData{}, false
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return FileData{}, false
	}

	return FileData{
		Id:        file.Id,
		Name:      file.Name,
		UpdatedAt: file.UpdatedAt,
		FileType:  file.FileType,
		Extension: file.Extension,
	}, true
}

func (repo *Repo) GetAllFiles(
	userId uuid.UUID,
	categId uuid.UUID,
) ([]EntityData, bool) {
	user, ok := repo.Users[userId]
	if !ok {
		return []EntityData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return []EntityData{}, false
	}

	var files []EntityData
	for _, file := range categ.Files {
		files = append(files, EntityData{
			Id:        file.Id,
			Name:      file.Name,
			UpdatedAt: file.UpdatedAt,
		})
	}

	return files, true
}

func (repo *Repo) GetFileAttachment(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) (AttachmentData, bool) {
	user, ok := repo.Users[userId]
	if !ok {
		return AttachmentData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return AttachmentData{}, false
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return AttachmentData{}, false
	}

	return AttachmentData{
		Id:        file.Id,
		Name:      file.Name,
		UpdatedAt: file.UpdatedAt,
		FileType:  file.FileType,
		Path:      file.Path,
	}, true
}

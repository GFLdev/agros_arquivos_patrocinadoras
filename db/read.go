package db

import "github.com/google/uuid"

func (repo *Repo) GetUserById(userId uuid.UUID) (QueryData, bool) {
	user, ok := repo.Users[userId]
	if repo.Users == nil || !ok {
		return QueryData{}, false
	}
	return QueryData{
		Id:        user.Id,
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}, true
}

func (repo *Repo) GetAllUsers() ([]QueryData, bool) {
	if repo.Users == nil {
		return []QueryData{}, false
	}

	var users []QueryData
	for _, user := range repo.Users {
		users = append(users, QueryData{
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
) (QueryData, bool) {
	user, ok := repo.Users[userId]
	if repo.Users == nil || !ok || user.Categories == nil {
		return QueryData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok {
		return QueryData{}, false
	}
	return QueryData{categ.Id, categ.Name, categ.UpdatedAt}, true
}

func (repo *Repo) GetAllCategories(userId uuid.UUID) ([]QueryData, bool) {
	user, ok := repo.Users[userId]
	if repo.Users == nil || !ok || user.Categories == nil {
		return []QueryData{}, false
	}

	categs := make([]QueryData, len(user.Categories))
	for _, categ := range user.Categories {
		categs = append(categs, QueryData{categ.Id, categ.Name, categ.UpdatedAt})
	}
	return categs, true
}

func (repo *Repo) GetFileById(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) (QueryData, bool) {
	user, ok := repo.Users[userId]
	if repo.Users == nil || !ok || user.Categories == nil {
		return QueryData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok || categ.Files == nil {
		return QueryData{}, false
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return QueryData{}, false
	}
	return QueryData{file.Id, file.Name, file.UpdatedAt}, true
}

func (repo *Repo) GetAllFiles(
	userId uuid.UUID,
	categId uuid.UUID,
) ([]QueryData, bool) {
	user, ok := repo.Users[userId]
	if repo.Users == nil || !ok || user.Categories == nil {
		return []QueryData{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok || categ.Files == nil {
		return []QueryData{}, false
	}

	files := make([]QueryData, len(categ.Files))
	for _, file := range categ.Files {
		files = append(files, QueryData{file.Id, file.Name, file.UpdatedAt})
	}
	return files, true
}

func (repo *Repo) GetFileAttachment(
	userId uuid.UUID,
	categId uuid.UUID,
	fileId uuid.UUID,
) (FileAttachment, bool) {
	user, ok := repo.Users[userId]
	if repo.Users == nil || !ok || user.Categories == nil {
		return FileAttachment{}, false
	}

	categ, ok := user.Categories[categId]
	if !ok || categ.Files == nil {
		return FileAttachment{}, false
	}

	file, ok := categ.Files[fileId]
	if !ok {
		return FileAttachment{}, false
	}
	return FileAttachment{
		Id:        file.Id,
		Name:      file.Name,
		UpdatedAt: file.UpdatedAt,
		FileType:  file.FileType,
		Path:      file.Path,
	}, true
}

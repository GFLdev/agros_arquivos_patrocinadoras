package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// DownloadHandler gerencia o envio de um arquivo para download.
func DownloadHandler(c echo.Context) error {
	ctx := GetAppContext(c)

	if err := checkAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	ctx.Repo.Lock()
	attach, ok := ctx.Repo.GetFileAttachment(userId, categId, fileId)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	defer ctx.Repo.Unlock()

	return c.Attachment(attach.Path, attach.Name)
}

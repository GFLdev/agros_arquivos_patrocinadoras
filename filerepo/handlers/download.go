package handlers

import (
	"agros_arquivos_patrocinadoras/filerepo/services"
	"agros_arquivos_patrocinadoras/filerepo/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// DownloadHandler gerencia o envio de um arquivo para download.
func DownloadHandler(c echo.Context) error {
	ctx := services.GetContext(c)

	if err := utils.CheckAuthentication(c); err != nil {
		return err
	}

	userId := uuid.MustParse(c.Param("userId"))
	categId := uuid.MustParse(c.Param("categId"))
	fileId := uuid.MustParse(c.Param("fileId"))

	ctx.FSServ.Mux.Lock()
	attach, ok := ctx.FSServ.FS.GetFileAttachment(userId, categId, fileId)
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}
	defer ctx.FSServ.Mux.Unlock()

	return c.Attachment(attach.Path, attach.Name)
}

package services

import (
	"agros_arquivos_patrocinadoras/filerepo/services/config"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"database/sql"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AppWrapper struct {
	Logger     *zap.Logger
	Config     *config.Config
	FileSystem *fs.FileSystem
	DB         *sql.DB
}

// GetContext recupera o AppWrapper do contexto do echo.
func GetContext(c echo.Context) *AppWrapper {
	return c.Get("appContext").(*AppWrapper)
}

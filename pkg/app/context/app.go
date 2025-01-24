package context

import (
	"agros_arquivos_patrocinadoras/pkg/app/fs"
	"agros_arquivos_patrocinadoras/pkg/types/config"
	"database/sql"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Context struct {
	Logger     *zap.Logger
	Config     *config.Config
	FileSystem *fs.FileSystem
	DB         *sql.DB
}

// GetContext recupera o AppContext do contexto do echo.
func GetContext(c echo.Context) *Context {
	return c.Get("appContext").(*Context)
}

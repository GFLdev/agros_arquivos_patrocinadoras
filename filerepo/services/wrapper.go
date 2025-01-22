package services

import (
	"agros_arquivos_patrocinadoras/filerepo/services/config"
	"agros_arquivos_patrocinadoras/filerepo/services/fs"
	"database/sql"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"sync"
)

type FSServ struct {
	FS  *fs.FS
	Mux *sync.Mutex
}

type AppWrapper struct {
	Logger *zap.Logger
	Config *config.Config
	FSServ *FSServ
	DB     *sql.DB
}

func NewAppWrapper(
	logr *zap.Logger,
	cfg *config.Config,
	fs *fs.FS,
	db *sql.DB,
) *AppWrapper {
	fsServ := &FSServ{
		FS:  fs,
		Mux: new(sync.Mutex),
	}

	return &AppWrapper{
		Logger: logr,
		Config: cfg,
		FSServ: fsServ,
		DB:     db,
	}
}

// GetContext recupera o AppWrapper do contexto do echo.
func GetContext(c echo.Context) *AppWrapper {
	return c.Get("appContext").(*AppWrapper)
}

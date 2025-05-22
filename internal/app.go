package app

import (
	"fmt"
	"os"
	"path/filepath"
	C "shortlink2/internal/cfg"
	D "shortlink2/internal/db"
	H "shortlink2/internal/http"
	L "shortlink2/internal/log"
	S "shortlink2/internal/service"
	T "shortlink2/internal/types"
	"time"
)

type App struct {
	hsrv T.IHTTPServer
	db   T.IDB
	log  T.ILog
	file string
}

func NewApp() *App {
	dir, file := execPathAndFname()
	cfg := C.NewCfgEnvMap(dir, file).Parse()
	log := L.NewLogFprintf(cfg, 3*time.Second, 4*time.Second)
	db := D.NewDBsqlite(cfg, log, dir)
	// db := D.NewDBmock(cfg, log)
	svcsl2 := S.NewSvcShortLink2(db, log)
	hsrv := H.NewHTTPServerNet(svcsl2, log, cfg)
	return &App{
		hsrv: hsrv,
		db:   db,
		log:  log,
		file: file,
	}
}

func (a *App) Start() func(err error) {
	logStop := a.log.Start()
	dbShutdown := a.db.ConnectDB()
	hsrvShutdown := a.hsrv.Run()
	a.log.LogInfo(a.file + " app started")
	return func(err error) {
		hsrvShutdown(err)
		dbShutdown(err)
		if err != nil {
			a.log.LogPanic(fmt.Errorf("%s: %w", a.file+" app stoped with error", err))
		} else {
			a.log.LogInfo(a.file + " app stoped")
		}
		logStop()
	}
}

func execPathAndFname() (string, string) {
	path, _ := os.Executable()
	return filepath.Split(path)
}

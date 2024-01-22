package app

import (
	"os"
	C "shortlink2/internal/cfg"
	H "shortlink2/internal/http"
	L "shortlink2/internal/log"
	T "shortlink2/internal/types"
	// D "shortlink2/internal/db"
	// S "shortlink2/internal/service"
)

type App struct {
	hsrv T.IHTTPServer
	db   T.IDB
	log  T.ILog
	cfg  *T.CfgEnv
}

func NewApp() *App {
	cfg := C.NewCfg(os.Args[0] + ".env")
	log := L.NewLogFprintf(cfg)
	db := D.NewSQLite(cfg, log)
	svcsl := S.NewSvcShortLink(db, log)
	hsrv := H.NewHTTPServerNet(svcsl, log, cfg)
	return &App{
		hsrv: hsrv,
		db:   db,
		log:  log,
		cfg:  cfg,
	}
}

func (a *App) Start() func(err error) {
	dbShutdown := a.db.ConnectDB()
	hsrvShutdown := a.hsrv.Run()
	a.log.LogInfo(a.cfg.SL_APP_NAME + " app started")
	return func(err error) {
		hsrvShutdown(err)
		dbShutdown(err)
		if err != nil {
			a.log.LogError(err, a.cfg.SL_APP_NAME+" app stoped with error")
		} else {
			a.log.LogInfo(a.cfg.SL_APP_NAME + " app stoped")
		}
	}
}

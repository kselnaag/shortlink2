package app

import (
	C "shortlink2/internal/cfg"
	D "shortlink2/internal/db"
	H "shortlink2/internal/http"
	L "shortlink2/internal/log"
	S "shortlink2/internal/service"
	T "shortlink2/internal/types"
)

type App struct {
	hsrv T.IHTTPServer
	db   T.IDB
	log  T.ILog
	cfg  *T.CfgEnv
}

func NewApp(cfgfilename string) *App {
	cfg := C.NewCfgEnv(cfgfilename)
	log := L.NewLogFprintf(cfg)
	// db := D.NewSQLite(cfg, log)
	db := D.NewDBMock(cfg, log)
	svcsl2 := S.NewSvcShortLink2(db, log)
	hsrv := H.NewHTTPServerNet(svcsl2, log, cfg)
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

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	R "shortlink2/internal/http/route"
	T "shortlink2/internal/types"
	W "shortlink2/web"
)

var _ T.IHTTPServer = (*HTTPServerNet)(nil)

type HTTPMess struct {
	Method string `json:"M"`
	Hash   string `json:"H"`
	Link   string `json:"L"`
}

type HTTPServerNet struct {
	hsrv *http.Server
	svc  T.ISvcShortLink2
	log  T.ILog
	cfg  *T.CfgEnv
	fs   http.FileSystem
}

func NewHTTPServerNet(svc T.ISvcShortLink2, log T.ILog, cfg *T.CfgEnv) *HTTPServerNet {
	subFS, err := fs.Sub(W.StaticFS, "data")
	if err != nil {
		log.LogError(err, "staticFS: embedFS error")
	}
	return &HTTPServerNet{
		hsrv: nil,
		svc:  svc,
		log:  log,
		cfg:  cfg,
		fs:   http.FS(subFS),
	}
}

/*
	curl -i -X POST localhost:8080/save -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"IsResp":false,"Mode":"client","Body":"http://lib.ru/PROZA/"}'
	Cache-Control: no-cache | Content-Type: text/html; charset=utf-8
	(5clp60)http://lib.ru (dhiu79)http://google.ru (8b4s29)http://lib.ru/PROZA/
*/

func (hns *HTTPServerNet) getRedirect(w http.ResponseWriter, r *http.Request) {
	hash, _ := strings.CutPrefix(r.URL.Path, "/")
	link := hns.svc.GetLinkLongFromLinkShort(hash)
	if len(link) == 0 {
		http.Error(w, "", http.StatusNotFound)
	} else {
		http.Redirect(w, r, link, http.StatusFound)
	}
}

func (hns *HTTPServerNet) postLoad(w http.ResponseWriter, r *http.Request) {
	mess := HTTPMess{}
	err := json.NewDecoder(r.Body).Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	link := hns.svc.GetLinkLongFromLinkShort(mess.Hash)
	if len(link) == 0 {
		http.Error(w, "", http.StatusNotFound)
	} else {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"M":"200","H":"%s","L":"%s"}`+"\n", mess.Hash, link)
	}
}

func (hns *HTTPServerNet) postSave(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "save %s\n", r.URL.Path)
	fmt.Println("save")
}

func (hns *HTTPServerNet) postDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "delete %s\n", r.URL.Path)
	fmt.Println("delete")
}

func (hns *HTTPServerNet) handlers() *R.RouteHandler {
	middlewares := &[]*R.Middleware{
		// R.NewMiddleware(midf1),
		// R.NewMiddleware(midf2),
	}
	routes := []*R.Route{
		R.NewRoute("GET", "/[a-z0-9]{6}", hns.getRedirect),
		R.NewRoute("POST", "/load", hns.postLoad),
		R.NewRoute("POST", "/save", hns.postSave),
		R.NewRoute("POST", "/delete", hns.postDelete),
	}
	staticfs := http.StripPrefix("/", http.FileServer(hns.fs))
	return R.NewRouteHandler(middlewares, routes, staticfs, hns.log)
}

func (hns *HTTPServerNet) Run() func(e error) {
	hns.hsrv = &http.Server{
		Addr:           hns.cfg.SL_HTTP_PORT,
		Handler:        hns.handlers(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				hns.log.LogPanic(err.(error), "Run(): net/http server panic")
			}
		}()
		err := hns.hsrv.ListenAndServe()
		if (err != nil) && (err != http.ErrServerClosed) {
			hns.log.LogError(err, "Run(): net/http server closed with error")
			os.Exit(1)
		}
		if err == http.ErrServerClosed {
			hns.log.LogInfo("net/http server closed")
		}
	}()
	hns.log.LogInfo("net/http server opened")
	return func(e error) {
		ctxSHD, cancelSHD := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelSHD()
		if err := hns.hsrv.Shutdown(ctxSHD); err != nil {
			hns.log.LogError(err, "Run(): net/http server graceful_shutdown error")
		}
		if e != nil {
			hns.log.LogError(e, "Run(): net/http server shutdown with error")
		}
	}
}

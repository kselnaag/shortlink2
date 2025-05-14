package route

import (
	"fmt"
	"net/http"
	"regexp"
	T "shortlink2/internal/types"
)

type Route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func NewRoute(method, pattern string, handler http.HandlerFunc) *Route {
	return &Route{
		method,
		regexp.MustCompile("^" + pattern + "$"),
		handler,
	}
}

type Middleware struct {
	handler http.HandlerFunc
}

func NewMiddleware(handler http.HandlerFunc) *Middleware {
	return &Middleware{
		handler,
	}
}

type RouteHandler struct {
	middlewares *[]*Middleware
	routes      []*Route
	staticfs    http.Handler
	log         T.ILog
}

func NewRouteHandler(middlewares *[]*Middleware, routes []*Route, staticfs http.Handler, log T.ILog) *RouteHandler {
	return &RouteHandler{
		middlewares,
		routes,
		staticfs,
		log,
	}
}

func (rh *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rh.log.LogError(fmt.Errorf("%s: %w", "500: some handler panics", err.(error)))
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}
	}()

	if rh.middlewares != nil {
		for _, middlwr := range *rh.middlewares {
			middlwr.handler(w, r)
		}
	}

	if len(rh.routes) == 0 {
		rh.log.LogError(fmt.Errorf("%s: %s", "500 internal server error", "empty routes"))
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	isWrongMethod := false
	for _, route := range rh.routes {
		if route.regex.MatchString(r.URL.Path) {
			if r.Method != route.method {
				isWrongMethod = true
				continue
			}
			route.handler(w, r)
			return
		}
	}
	if isWrongMethod {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if rh.staticfs != nil {
		rh.staticfs.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

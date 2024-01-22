package route

import (
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
	log         T.ILog
}

func NewRouteHandler(middlewares *[]*Middleware, routes []*Route, log T.ILog) *RouteHandler {
	return &RouteHandler{
		middlewares,
		routes,
		log,
	}
}

func (rh *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rh.log.LogError(err.(error), "some handler panics")
		}
	}()
	if rh.middlewares != nil {
		for _, middlwr := range *rh.middlewares {
			middlwr.handler(w, r)
		}
	}
	if rh.routes == nil {
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
	http.NotFound(w, r)
}
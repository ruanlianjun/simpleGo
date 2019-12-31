package simple

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(c *Context)

type (
	RouterGroup struct {
		prefix string
		//middlewares []HandleFunc
		parent      *RouterGroup
		simple      *Simple
		middlewares []MiddlewareHandle
	}
	Simple struct {
		router *router
		*RouterGroup
		groups []*RouterGroup
	}
)

func New() *Simple {
	Simple := &Simple{router: NewRouter()}
	Simple.RouterGroup = &RouterGroup{simple: Simple}
	Simple.groups = []*RouterGroup{Simple.RouterGroup}

	return Simple
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	Simple := group.simple
	newGroup := &RouterGroup{
		prefix: group.prefix + "/" + prefix,
		parent: group,
		simple: Simple,
	}
	Simple.groups = append(Simple.groups, newGroup)

	return newGroup
}

//func (group *RouterGroup) Use(middlewares ...HandleFunc) {
//	group.middlewares = append(group.middlewares, middlewares...)
//}

func (group *RouterGroup) addRoute(method string,
	comp string,
	handle HandleFunc) string {
	pattern := group.prefix + "/" + comp
	group.simple.router.addRoute(method, pattern, handle)
	return pattern
}

func (group *RouterGroup) GET(patter string, handle HandleFunc) *RouterGroup {
	group.addRoute("GET", patter, handle)
	return group
}

func (group *RouterGroup) POST(patter string, handle HandleFunc) *RouterGroup {
	group.addRoute("POST", patter, handle)
	return group
}

func (group *RouterGroup) Middleware(middleware MiddlewareHandle) *RouterGroup {
	group.middlewares = append(group.middlewares, middleware)

	return group
}

func (h *Simple) Run(addr string) (error error) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println(h.prefix)

	return http.ListenAndServe(addr, h)
}

func (h *Simple) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []MiddlewareHandle

	for _, group := range h.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			fmt.Printf("host=%s middleware=%v group_prefix=%s\n",
				r.URL.Path, group.middlewares, group.prefix)
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	context := NewContext(w, r)
	context.middlewares = middlewares

	h.router.handle(context)
}

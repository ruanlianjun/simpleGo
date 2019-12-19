package simple

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type HandleFunc func(c *Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandleFunc
		parent      *RouterGroup
		simple      *Simple
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

func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string,
	comp string,
	handle HandleFunc) {
	pattern := group.prefix + "/" + comp
	group.simple.router.addRoute(method, pattern, handle)

}

func (group *RouterGroup) GET(patter string, handle HandleFunc) {
	//fmt.Println(group)
	group.addRoute("GET", patter, handle)
}
func (group *RouterGroup) POST(patter string, handle HandleFunc) {
	group.addRoute("POST", patter, handle)
}

func (h *Simple) Run(addr string) (error error) {

	now := time.Now().String()
	log.Println("run now" + now + " at" + addr)

	return http.ListenAndServe(addr, h)
}

func (h *Simple) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandleFunc

	for _, group := range h.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	context := NewContext(w, r)
	context.handlers = middlewares
	h.router.handle(context)
}

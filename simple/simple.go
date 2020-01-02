package simple

import (
	"net/http"
	"strings"
)

type RouterGroup struct {
	prefix      string
	engine      *Engine
	parent      *RouterGroup
	middlewares []MiddlewareHandle
}

type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

type HandleFunc func(c *Context)
type MiddlewareHandle func(c *Context, next func(c *Context))

func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine

	prefix = strings.Trim(prefix, "/")
	newGroup := &RouterGroup{
		prefix: group.prefix + "/" + prefix,
		engine: engine,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}


func (group *RouterGroup) addRoute(method string,
	comp string,
	handle HandleFunc) {
	pattern := group.prefix + comp

	group.engine.router.addRoute(method, pattern, handle)
}

func (group *RouterGroup) GET(patter string, handle HandleFunc) {
	group.addRoute("GET", patter, handle)
}

func (group *RouterGroup) POST(patter string, handle HandleFunc) {
	group.addRoute("POST", patter, handle)
}

func (group *RouterGroup) Middleware(middleware MiddlewareHandle) *RouterGroup {
	group.middlewares = append(group.middlewares, middleware)
	return group
}

func (engine *Engine) Run(addr string) (error error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	var middlewares []MiddlewareHandle

	for _, group := range engine.groups {
		if strings.HasPrefix(rq.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	context := NewContext(w, rq)
	context.handlers = middlewares
	engine.router.handle(context)
}

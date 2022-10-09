package tinyweb

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Engine struct {
	router *router
	*RouterGroup
	groups   []*RouterGroup
	template *template.Template
	funcMap  template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine, prefix: ""}
	engine.Use(Recovery())
	return engine
}

func (ng *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)
	ctx.engine = ng
	handlerKey, params, ok := ng.router.getRouter(ctx.Method, ctx.Path)
	if ok {
		ctx.handlers = append(append([]HandlerFunc{}, ng.middlewares...), ctx.handlers...)
		for _, group := range ng.groups {
			if strings.HasPrefix(handlerKey, group.prefix) {
				ctx.handlers = append(ctx.handlers, group.middlewares...)
			}
		}

		ctx.Params = params
		key := ctx.Method + "-" + handlerKey
		ctx.handlers = append(ctx.handlers, ng.router.handlers[key])
	} else {
		ctx.handlers = append(ctx.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	ctx.Next()
}

func (ng *Engine) GET(pattern string, handler HandlerFunc) {
	ng.addRoute("GET", pattern, handler)
}

func (ng *Engine) POST(pattern string, handler HandlerFunc) {
	ng.addRoute("POST", pattern, handler)
}

func (ng *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	ng.router.addRoute(method, pattern, handler)
}

func (ng *Engine) Run(addr string) (err error) {
	log.Printf("serve at http://%s\n", addr)
	return http.ListenAndServe(addr, ng)
}

func (ng *Engine) SetFuncMap(funcMap template.FuncMap) {
	ng.funcMap = funcMap
}

func (ng *Engine) LoadHTMLGlob(pattern string) {
	ng.template = template.Must(template.New("").Funcs(ng.funcMap).ParseGlob(pattern))
}

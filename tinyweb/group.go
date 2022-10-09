package tinyweb

import (
	"net/http"
	"path"
)

type RouterGroup struct {
	prefix      string
	engine      *Engine
	middlewares []HandlerFunc
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	g := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: group.engine,
	}
	group.engine.groups = append(group.engine.groups, g)
	return g
}

func (group *RouterGroup) addRoute(method string, path string, handler HandlerFunc) {
	group.engine.router.addRoute(method, group.prefix+path, handler)
}

func (group *RouterGroup) getRoute(method string, path string) (string, map[string]string, bool) {
	return group.engine.router.getRouter(method, group.prefix+path)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) Static(relativePath string, root string, staticFiles []string) {
	fs := http.Dir(root)
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	for _, x := range staticFiles {
		urlPattern := path.Join(relativePath, x)
		group.GET(urlPattern, func(c *Context) {
			file := c.Params["filepath"]
			if _, err := fs.Open(file); err != nil {
				c.SetStatus(http.StatusNotFound)
				return
			}

			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	}
}

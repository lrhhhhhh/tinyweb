package tinyweb

import (
	"errors"
	"fmt"
)

type HandlerFunc func(*Context)

type router struct {
	handlers map[string]HandlerFunc
	path     map[string]*trie
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		path:     make(map[string]*trie),
	}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	_, ok := r.path[method]
	if !ok {
		r.path[method] = newTrie()
	}
	key, ok := r.path[method].insert(pattern)
	if ok {
		key = method + "-" + key
		r.handlers[key] = handler
	} else {
		panic(errors.New(fmt.Sprintf("invalid path: %s\n", pattern)))
	}
}

func (r *router) getRouter(method string, path string) (string, map[string]string, bool) {
	_, ok := r.path[method]
	if !ok {
		return "", nil, false
	}
	return r.path[method].find(path)
}

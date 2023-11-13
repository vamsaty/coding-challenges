package server

import "strings"

/* ---------------- Registry ---------------- */

// Registry stores the supported URLs and their handlers
type Registry struct {
	//  store represents map[Method]map[Path]Handler
	store map[string]map[string]MyHttpHandler
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (registry *Registry) Register(method, path string, Handler MyHttpHandler) {
	if registry.store == nil {
		registry.store = make(map[string]map[string]MyHttpHandler)
	}
	if registry.store[method] == nil {
		registry.store[method] = make(map[string]MyHttpHandler)
	}
	registry.store[method][path] = Handler
}

func (registry *Registry) GetHandler(method, path string) MyHttpHandler {
	// match prefix of store
	for pathKey, handler := range registry.store[method] {
		if strings.HasPrefix(path, pathKey) {
			return handler
		}
	}
	return HandlerNotFound
}

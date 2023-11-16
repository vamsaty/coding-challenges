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
	// NOTE: match prefix of store. If there are paths starting with the same
	// prefix, the first one is returned. if "/" is registered and "/api" is
	// requested, "/" is returned. This can lead to only "/" being returned
	for pathKey, handler := range registry.store[method] {
		if strings.HasPrefix(path, pathKey) {
			return handler
		}
	}
	return HandlerNotFound
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type Router struct {
	*mux.Router
	middleware   alice.Chain
	healthchecks map[string]HealthCheckHandler
}

// NewRouter constructor
func NewRouter() *Router {
	return &Router{
		Router:       mux.NewRouter(),
		healthchecks: make(map[string]HealthCheckHandler),
	}
}

// AddHealthCheck handler
func (r *Router) AddHealthCheck(name string, handle HealthCheckHandler) error {
	if _, ok := r.healthchecks[name]; ok {
		return fmt.Errorf("the %s healthcheck handler already exists", name)
	}

	r.healthchecks[name] = handle

	return nil
}

// EnableHealthCheck endpoint
func (r *Router) EnableHealthCheck() {
	r.AddRouteFunc("/health", r.healthHandler).Methods("GET", "HEAD")
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	code := http.StatusOK

	response := healthCheckProcessor(req.Context(), r.healthchecks)

	if !response.Status {
		code = http.StatusServiceUnavailable
	}

	JSON(w, code, response)
}

func (r *Router) Use(middleware ...alice.Constructor) {
	r.middleware = r.middleware.Append(middleware...)
}

func (r *Router) AddController(controller Controller) {
	controller.Mount(r)
}

func (r *Router) AddRoute(path string, handler http.Handler) *mux.Route {
	return r.Handle(path, r.middleware.Then(handler))
}

func (r *Router) AddRouteFunc(path string, handler http.HandlerFunc) *mux.Route {
	return r.Handle(path, r.middleware.ThenFunc(handler))
}

func (r *Router) AddPrefixRoute(prefix string, handler http.Handler) *mux.Route {
	return r.PathPrefix(prefix).Handler(r.middleware.Then(handler))
}

func (r *Router) AddPrefixRouteFunc(prefix string, handler http.HandlerFunc) *mux.Route {
	return r.PathPrefix(prefix).Handler(r.middleware.ThenFunc(handler))
}

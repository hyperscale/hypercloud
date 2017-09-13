// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/hyperscale/hyperpaas/server"
	"github.com/hyperscale/hyperpaas/services"
)

// APIController struct
type APIController struct {
	authenticationService *services.AuthenticationService
	userService           *services.UserService
}

// NewAPIController func
func NewAPIController(authenticationService *services.AuthenticationService, userService *services.UserService) (*APIController, error) {
	return &APIController{
		authenticationService: authenticationService,
		userService:           userService,
	}, nil
}

// Mount endpoints
func (c APIController) Mount(r *server.Router) {
	r.AddRouteFunc("/api/registries", c.GetRegistriesHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/api/registries", c.PostRegistryHandler).Methods(http.MethodPost)
	r.AddRouteFunc("/api/registries/{id:[0-9]+}", c.PutRegistryHandler).Methods(http.MethodPut)
	r.AddRouteFunc("/api/registries/{id:[0-9]+}", c.DeleteRegistryHandler).Methods(http.MethodDelete)
	r.AddRouteFunc("/api/registries/{id:[0-9]+}/repositories", c.GetRegistryRepositoriesHandler).Methods(http.MethodGet)

	r.AddRouteFunc("/api/stacks", c.GetStacksHandler).Methods(http.MethodGet)

	r.AddRouteFunc("/api/authenticate", c.PostAuthHandler).Methods(http.MethodPost)

	r.AddRouteFunc("/api/users", c.PostUserHandler).Methods(http.MethodPost)
	r.AddRouteFunc("/api/me", c.GetMeHandler).Methods(http.MethodGet)
}

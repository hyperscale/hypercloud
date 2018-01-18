// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/asdine/storm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/euskadi31/go-server"
	std "github.com/euskadi31/go-std"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hyperscale/hyperpaas/database/entity"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/rs/zerolog/log"
)

// StackController struct
type StackController struct {
	dockerClient *docker.Client
	db           *storm.DB
	validator    *server.Validator
}

// NewStackController func
func NewStackController(dockerClient *docker.Client, db *storm.DB, validator *server.Validator) (*StackController, error) {
	if err := db.Init(&entity.Stack{}); err != nil {
		return nil, err
	}

	return &StackController{
		dockerClient: dockerClient,
		db:           db,
		validator:    validator,
	}, nil
}

// Mount endpoints
func (c StackController) Mount(r *server.Router) {
	r.AddRouteFunc("/v1/stacks", c.getStacksHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/v1/stacks", c.postStacksHandler).Methods(http.MethodPost)
	// r.AddRouteFunc("/v1/stacks/{id:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}", c.getStackHandler).Methods(http.MethodGet)
	// r.AddRouteFunc("/v1/stacks/{id:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}/services", c.getStackServicesHandler).Methods(http.MethodGet)

	r.AddRouteFunc("/v1/stacks/{id:[a-z]+(?:-[a-z0-9]+)*}/services", c.getStackServicesHandler).Methods(http.MethodGet)
}

// swagger:route GET /v1/stacks Stack getStacksHandler
//
// Get the stacks list
//
//     Responses:
//       200: Stack
//
func (c StackController) getStacksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stackMap := map[string]*docker.Stack{}

	{
		var stacks []*entity.Stack

		if err := c.db.All(&stacks); err != nil {
			log.Error().Err(err).Msg("Get All Stacks")

			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		for _, stack := range stacks {
			stackMap[stack.Name] = &docker.Stack{
				Name:     stack.Name,
				Services: 0,
			}
		}
	}

	{
		stacks, err := c.dockerClient.StackList(ctx)
		if err != nil {
			log.Error().Err(err).Msg("StackList")

			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		for _, stack := range stacks {
			stackMap[stack.Name] = &stack
		}
	}

	response := []*docker.Stack{}

	for _, stack := range stackMap {
		response = append(response, stack)
	}

	server.JSON(w, http.StatusOK, response)
}

// swagger:route POST /v1/stacks Stack postStacksHandler
//
// Create stack
//
//     Responses:
//       201: Stack
//
func (c StackController) postStacksHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("uuid.NewRandom")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	stack := &entity.Stack{
		ID: id.String(),
	}

	if err := json.NewDecoder(r.Body).Decode(stack); err != nil {
		log.Error().Err(err).Msg("Decode body request")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	if result := c.validator.Validate("stack", stack); !result.IsValid() {
		log.Error().Err(result.AsError()).Msg("Validate body request")

		server.FailureFromValidator(w, result)

		return
	}

	stack.CreatedAt = std.DateTimeFrom(time.Now().UTC())

	if err := c.db.Save(stack); err != nil {
		log.Error().Err(err).Msg("Save Stack")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, stack)
}

// swagger:route GET /v1/stacks/{id} Stack getStackHandler
//
// Get a stack by id
//
//     Responses:
//       200: Stack
//
func (c StackController) getStackHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	log.Debug().Msgf("Stack ID: %s", id)

	stack := &entity.Stack{}

	if err := c.db.One("ID", id, stack); err != nil {
		server.FailureFromError(w, http.StatusNotFound, err)

		return
	}

	server.JSON(w, http.StatusOK, stack)
}

// swagger:route POST /v1/stacks/{id}/services Stack getStackServicesHandler
//
// Get stack services list
//
//     Responses:
//       200: Service
//
func (c StackController) getStackServicesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)

	id := params["id"]

	log.Debug().Msgf("Stack ID: %s", id)

	filter := filters.NewArgs()
	filter.Add("label", "com.docker.stack.namespace="+id)

	services, err := c.dockerClient.ServiceList(ctx, types.ServiceListOptions{
		Filters: filter,
	})
	if err != nil {
		log.Error().Err(err).Msg("ServiceList")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, services)
}

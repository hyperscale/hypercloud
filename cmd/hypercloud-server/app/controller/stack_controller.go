// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/asdine/storm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	server "github.com/euskadi31/go-server"
	"github.com/euskadi31/go-server/request"
	"github.com/euskadi31/go-server/response"
	"github.com/gorilla/mux"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/database/entity"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
	"github.com/rs/zerolog/log"
)

// StackController struct
type StackController struct {
	dockerClient *docker.Client
	db           *storm.DB
	validator    *request.Validator
}

// NewStackController func
func NewStackController(dockerClient *docker.Client, db *storm.DB, validator *request.Validator) (*StackController, error) {
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
	r.HandleFunc("/v1/stacks", c.getStacksHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/stacks", c.postStacksHandler).Methods(http.MethodPost)
	// r.HandleFunc("/v1/stacks/{id:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}", c.getStackHandler).Methods(http.MethodGet)
	// r.HandleFunc("/v1/stacks/{id:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}/services", c.getStackServicesHandler).Methods(http.MethodGet)

	r.HandleFunc("/v1/stacks/{id:[a-z]+(?:-[a-z0-9]+)*}/services", c.getStackServicesHandler).Methods(http.MethodGet)
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

	stackMap := map[string]docker.Stack{}

	{
		var stacks []*entity.Stack

		if err := c.db.All(&stacks); err != nil {
			log.Error().Err(err).Msg("Get All Stacks")

			response.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		for _, stack := range stacks {
			stackMap[stack.Name] = docker.Stack{
				Name:     stack.Name,
				Services: 0,
			}
		}
	}

	{
		stacks, err := c.dockerClient.StackList(ctx)
		if err != nil {
			log.Error().Err(err).Msg("StackList")

			response.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}

		for _, stack := range stacks {
			stackMap[stack.Name] = stack
		}
	}

	resp := []docker.Stack{}

	for _, stack := range stackMap {
		resp = append(resp, stack)
	}

	response.Encode(w, r, http.StatusOK, resp)
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

	stack := &entity.Stack{}

	if err := json.NewDecoder(r.Body).Decode(stack); err != nil {
		log.Error().Err(err).Msg("Decode body request")

		response.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	if result := c.validator.Validate("stack", stack); !result.IsValid() {
		log.Error().Err(result.AsError()).Msg("Validate body request")

		response.FailureFromValidator(w, result)

		return
	}

	if err := c.db.Save(stack); err != nil {
		log.Error().Err(err).Msg("Save Stack")

		response.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	response.Encode(w, r, http.StatusCreated, stack)
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
		response.FailureFromError(w, http.StatusNotFound, err)

		return
	}

	response.Encode(w, r, http.StatusOK, stack)
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

		response.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	response.Encode(w, r, http.StatusOK, services)
}

// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asdine/storm"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-std"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/hyperscale/hyperpaas/database/entity"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/http/request"
	"github.com/rs/zerolog/log"
)

// ServiceController struct
type ServiceController struct {
	dockerClient *docker.Client
	db           *storm.DB
	validator    *server.Validator
	queryDecoder *schema.Decoder
}

// NewServiceController func
func NewServiceController(dockerClient *docker.Client, db *storm.DB, validator *server.Validator) (*ServiceController, error) {
	if err := db.Init(&entity.Service{}); err != nil {
		return nil, err
	}

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	return &ServiceController{
		dockerClient: dockerClient,
		db:           db,
		validator:    validator,
		queryDecoder: decoder,
	}, nil
}

// Mount endpoints
func (c ServiceController) Mount(r *server.Router) {
	r.AddRouteFunc("/v1/services", c.getServicesHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/v1/services", c.postServiceHandler).Methods(http.MethodPost)
	r.AddRouteFunc("/v1/services/{id:[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}}", c.getServiceHandler).Methods(http.MethodGet)
}

// swagger:route GET /v1/services Service getServicesHandler
//
// Get the services list
//
//     Responses:
//       200: Service
//
func (c ServiceController) getServicesHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	query := &request.ServicesRequest{}

	if err := c.queryDecoder.Decode(query, r.URL.Query()); err != nil {
		log.Error().Err(err).Msg("Decode query parameters")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	var services []*entity.Service

	if query.StackID != "" {
		if err := c.db.Find("StackID", query.StackID, &services); err != nil {
			log.Error().Err(err).Msgf("Get All Services by StackID: %s", query.StackID)

			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}
	} else {
		if err := c.db.All(&services); err != nil {
			log.Error().Err(err).Msg("Get All Services")

			server.FailureFromError(w, http.StatusInternalServerError, err)

			return
		}
	}

	server.JSON(w, http.StatusOK, services)
}

// swagger:route POST /v1/services Service postServiceHandler
//
// Create service
//
//     Responses:
//       201: Service
//
func (c ServiceController) postServiceHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("uuid.NewRandom")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	service := &entity.Service{
		ID: id.String(),
	}

	if err := json.NewDecoder(r.Body).Decode(service); err != nil {
		log.Error().Err(err).Msg("Decode body request")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	if result := c.validator.Validate("service", service); !result.IsValid() {
		log.Error().Err(result.AsError()).Msg("Validate body request")

		server.FailureFromValidator(w, result)

		return
	}

	service.Hosts = append(service.Hosts, fmt.Sprintf("%s.%s", service.Name, "hyperpaas.service"))

	service.CreatedAt = std.DateTimeFrom(time.Now().UTC())

	if err := c.db.Save(service); err != nil {
		log.Error().Err(err).Msg("Save Service")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, service)
}

// swagger:route GET /v1/services/{id} Stack getServiceHandler
//
// Get a service by id
//
//     Responses:
//       200: Service
//
func (c ServiceController) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	log.Debug().Msgf("Service ID: %s", id)

	service := &entity.Service{}

	if err := c.db.One("ID", id, service); err != nil {
		server.FailureFromError(w, http.StatusNotFound, err)

		return
	}

	server.JSON(w, http.StatusOK, service)
}

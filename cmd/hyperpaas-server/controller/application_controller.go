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
	"github.com/docker/docker/client"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-std"
	"github.com/google/uuid"
	"github.com/hyperscale/hyperpaas/database/entity"
	"github.com/rs/zerolog/log"
)

// ApplicationController struct
type ApplicationController struct {
	dockerClient *client.Client
	db           *storm.DB
	validator    *server.Validator
}

// NewApplicationController func
func NewApplicationController(dockerClient *client.Client, db *storm.DB, validator *server.Validator) (*ApplicationController, error) {
	if err := db.Init(&entity.Application{}); err != nil {
		return nil, err
	}

	return &ApplicationController{
		dockerClient: dockerClient,
		db:           db,
		validator:    validator,
	}, nil
}

// Mount endpoints
func (c ApplicationController) Mount(r *server.Router) {
	r.AddRouteFunc("/v1/applications", c.GetApplicationsHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/v1/applications", c.PostApplicationGandler).Methods(http.MethodPost)
}

// GetApplicationsHandler endpoint
func (c ApplicationController) GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	var applications []*entity.Application

	if err := c.db.All(&applications); err != nil {
		log.Error().Err(err).Msg("Get All Applications")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, applications)
}

// PostApplicationGandler endpoint
func (c ApplicationController) PostApplicationGandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error().Err(err).Msg("uuid.NewRandom")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	application := &entity.Application{
		ID: id.String(),
	}

	if err := json.NewDecoder(r.Body).Decode(application); err != nil {
		log.Error().Err(err).Msg("Decode body request")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	if result := c.validator.Validate("application", application); !result.IsValid() {
		log.Error().Err(result.AsError()).Msg("Validate body request")

		server.FailureFromValidator(w, result)

		return
	}

	application.Hosts = append(application.Hosts, fmt.Sprintf("%s.%s", application.Slug(), "hyperpaas.service"))

	application.CreatedAt = std.DateTimeFrom(time.Now().UTC())

	if err := c.db.Save(application); err != nil {
		log.Error().Err(err).Msg("Save Application")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, application)
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hyperscale/hyperpaas/database"
	"github.com/hyperscale/hyperpaas/entity"
	"github.com/hyperscale/hyperpaas/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetRegistriesHandler /api/registries
func (c APIController) GetRegistriesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.FromContext(r.Context())
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	var registries []entity.Registry
	if err := db.All(&registries); err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, registries)
}

// PostRegistryHandler /api/registries
func (c APIController) PostRegistryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var registry entity.Registry

	if err := json.NewDecoder(r.Body).Decode(&registry); err != nil {
		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}
	defer r.Body.Close()

	/*dc := DockerFromContext(ctx)

	auth, err := dc.RegistryLogin(ctx, types.AuthConfig{
		Username:      registry.Username,
		Password:      registry.Password,
		ServerAddress: "https://" + registry.Server + "/v2/",
	})
	if err != nil {
		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	xlog.Debugf("Auth Registry: %#v", auth)
	*/
	db, err := database.FromContext(ctx)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	if err := db.Save(&registry); err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, registry)
}

// PutRegistryHandler /api/registries/{id:[0-9]+}
func (c APIController) PutRegistryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Info().Msgf("ID:", vars["id"])
}

// DeleteRegistryHandler /api/registries/{id:[0-9]+}
func (c APIController) DeleteRegistryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Info().Msgf("ID:", vars["id"])
}

// GetRegistryRepositoriesHandler /api/registries/{id:[0-9]+}/repositories
func (c APIController) GetRegistryRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		server.NotFoundFailure(w, r)

		return
	}

	db, err := database.FromContext(r.Context())
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	var registry entity.Registry

	if err := db.One("ID", ID, &registry); err != nil {
		server.FailureFromError(w, http.StatusNotFound, errors.Wrapf(err, "Cannot find registry by ID: %d", ID))

		return
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v2/_catalog", registry.Server), nil)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	// Add header with json of username and password
	req.SetBasicAuth(registry.Username, registry.Password)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	log.Debug().Msgf("Response: %s", string(b))

	//json.NewDecoder(req.Body).Decode(&)

	log.Info().Msgf("ID:", vars["id"])
}

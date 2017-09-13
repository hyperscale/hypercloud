// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/hyperscale/hyperpaas/entity"
	"github.com/hyperscale/hyperpaas/server"
	"github.com/rs/zerolog/log"
)

// PostAuthHandler /api/authenticate
func (c APIController) PostAuthHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	req := &entity.Auth{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	defer r.Body.Close()

	response, err := c.authenticationService.Authenticate(req)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusUnauthorized, err)
	}

	server.JSON(w, http.StatusOK, response)
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SermoDigital/jose/jws"
	"github.com/hyperscale/hyperpaas/entity"
	"github.com/hyperscale/hyperpaas/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetMeHandler /api/me
func (c APIController) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)

		ID, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Error().Err(err).Msg("")

			server.NotFoundFailure(w, r)

			return
		}
	*/
	token, err := jws.ParseFromHeader(r, jws.Compact)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	claims := token.Payload().(map[string]interface{})

	id, ok := claims["sub"]
	if !ok {
		err := errors.New("Missing subject in JWT payload")

		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	ID, err := strconv.Atoi(id.(string))
	if err != nil {
		log.Error().Err(err).Msg("")

		server.NotFoundFailure(w, r)

		return
	}

	user, err := c.userService.GetByID(ID)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusOK, user)
}

// PostUserHandler /api/users
func (c APIController) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	req := &entity.User{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	defer r.Body.Close()

	user, err := c.userService.Create(req)
	if err != nil {
		log.Error().Err(err).Msg("")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	server.JSON(w, http.StatusCreated, user)
}

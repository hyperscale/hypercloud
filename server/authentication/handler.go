// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package authentication

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hyperscale/hyperpaas/server"
	"github.com/rs/zerolog/log"
)

// NewAuthHandler middleware
func NewAuthHandler(provider Provider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Skip public endpoints
			if r.URL.Path == "/health" || r.URL.Path == "/api/authenticate" || (r.URL.Path == "/api/users" && r.Method == http.MethodPost) {
				next.ServeHTTP(w, r)

				return
			}

			authenticate := `Bearer realm="HyperPaaS"`

			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				authenticate += `, error="invalid_request", error_description="Authorization header invalid"`

				w.Header().Set("WWW-Authenticate", authenticate)

				log.Error().Msg("Authorization header invalid")

				server.FailureFromError(w, http.StatusBadRequest, errors.New("Bad Request"))

				return
			}

			token := s[1]

			log.Debug().Str("access_token", token).Msg("")

			if !provider.Validate(token) {
				authenticate += `, error="invalid_token", error_description="Access token invalid or expired"`

				w.Header().Set("WWW-Authenticate", authenticate)

				log.Error().Msg("Access token invalid or expired")

				server.FailureFromError(w, http.StatusUnauthorized, errors.New("Unauthorized"))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

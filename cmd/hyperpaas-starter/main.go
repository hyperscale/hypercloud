// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpaas/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = container.Get(ServiceLoggerKey).(zerolog.Logger)
	cfg := container.Get(ServiceConfigKey).(*config.Configuration)

	router := container.Get(ServiceRouterKey).(*server.Router)

	addr := cfg.Server.Addr()

	log.Info().Msgf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}

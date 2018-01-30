// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpaas/config"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = container.Get(ServiceLoggerKey).(zerolog.Logger)
	cfg := container.Get(ServiceConfigKey).(*config.Configuration)
	router := container.Get(ServiceRouterKey).(*server.Router)
	docker := container.Get(ServiceDockerKey).(*docker.Client)

	info, err := docker.Info(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Get Docker Info")

		return
	}

	if info.Swarm.LocalNodeState != "active" {
		log.Error().Msgf("Swarm Status: %s", info.Swarm.LocalNodeState)

		return
	}

	log.Info().Msgf("Swarm Status: %s", info.Swarm.LocalNodeState)

	if !info.Swarm.ControlAvailable {
		log.Error().Msgf("This node is not a Docker Swarm Manager")

		return
	}

	addr := cfg.Server.Addr()

	log.Info().Msgf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}

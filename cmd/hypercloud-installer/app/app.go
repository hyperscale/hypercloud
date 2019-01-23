// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	server "github.com/euskadi31/go-server"
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-installer/app/container"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Run hypercloud starter
func Run() (err error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	_ = service.Get(container.LoggerKey)

	docker := service.Get(container.DockerKey).(*docker.Client)

	info, err := docker.Info(context.Background())
	if err != nil {
		return errors.Wrap(err, "Get Docker Info")
	}

	if info.Swarm.LocalNodeState != "active" {
		return fmt.Errorf("Swarm Status: %s", info.Swarm.LocalNodeState)
	}

	log.Info().Msgf("Swarm Status: %s", info.Swarm.LocalNodeState)

	if !info.Swarm.ControlAvailable {
		return errors.New("This node is not a Docker Swarm Manager")
	}

	router := service.Get(container.RouterKey).(*server.Server)

	log.Info().Msg("Rinning")

	go func() {
		log.Info().Msg("Rinning HTTP Router")
		if e := router.Run(); e != nil {
			err = errors.Wrap(e, "server.Run")
		}
	}()

	<-sig

	log.Info().Msg("Shutdown")

	return router.Shutdown()
}

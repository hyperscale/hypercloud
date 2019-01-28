// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"github.com/asdine/storm"
	"github.com/euskadi31/go-server/request"
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-server/app/config"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-server/app/controller"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
	"github.com/rs/zerolog/log"
)

// Services keys
const (
	TopologyControllerKey = "service.controller.topology"
	ServiceControllerKey  = "service.controller.service"
	DeployControllerKey   = "service.controller.deploy"
	StackControllerKey    = "service.controller.stack"
	VersionControllerKey  = "service.controller.version"
	DockerControllerKey   = "service.controller.docker"
	EventControllerKey    = "service.controller.event"
)

func init() {
	service.Set(TopologyControllerKey, func(c service.Container) interface{} {
		dockerClient := c.Get(DockerKey).(*docker.Client)

		controller, err := controller.NewTopologyController(dockerClient)
		if err != nil {
			log.Fatal().Err(err).Msg(TopologyControllerKey)
		}

		return controller
	})

	service.Set(ServiceControllerKey, func(c service.Container) interface{} {
		dockerClient := c.Get(DockerKey).(*docker.Client)
		db := c.Get(DBKey).(*storm.DB)
		validator := c.Get(ValidatorKey).(*request.Validator)

		controller, err := controller.NewServiceController(dockerClient, db, validator)
		if err != nil {
			log.Fatal().Err(err).Msg(ServiceControllerKey)
		}

		return controller
	})

	service.Set(DeployControllerKey, func(c service.Container) interface{} {
		dockerClient := c.Get(DockerKey).(*docker.Client)
		db := c.Get(DBKey).(*storm.DB)
		validator := c.Get(ValidatorKey).(*request.Validator)

		controller, err := controller.NewDeployController(dockerClient, db, validator)
		if err != nil {
			log.Fatal().Err(err).Msg(DeployControllerKey)
		}

		return controller
	})

	service.Set(StackControllerKey, func(c service.Container) interface{} {
		dockerClient := c.Get(DockerKey).(*docker.Client)
		db := c.Get(DBKey).(*storm.DB)
		validator := c.Get(ValidatorKey).(*request.Validator)

		controller, err := controller.NewStackController(dockerClient, db, validator)
		if err != nil {
			log.Fatal().Err(err).Msg(StackControllerKey)
		}

		return controller
	})

	service.Set(VersionControllerKey, func(c service.Container) interface{} {
		controller, err := controller.NewVersionController()
		if err != nil {
			log.Fatal().Err(err).Msg(VersionControllerKey)
		}

		return controller
	})

	service.Set(DockerControllerKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)

		controller, err := controller.NewDockerController(cfg.Docker.Host)
		if err != nil {
			log.Fatal().Err(err).Msg(DockerControllerKey)
		}

		return controller
	})

	service.Set(EventControllerKey, func(c service.Container) interface{} {
		dockerClient := c.Get(DockerKey).(*docker.Client)

		controller, err := controller.NewEventController(dockerClient)
		if err != nil {
			log.Fatal().Err(err).Msg(EventControllerKey)
		}

		return controller
	})
}

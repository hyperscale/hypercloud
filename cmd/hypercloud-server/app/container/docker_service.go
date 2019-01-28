// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"os"

	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
	"github.com/rs/zerolog/log"
)

// Services keys
const (
	DockerKey = "service.docker"
)

func init() {
	service.Set(DockerKey, func(c service.Container) interface{} {
		os.Setenv("DOCKER_API_VERSION", "1.39")

		dc, err := docker.NewEnvClient()
		if err != nil {
			log.Fatal().Err(err).Msg(DockerKey)
		}

		return dc // *docker.Client
	})
}

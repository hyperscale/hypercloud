// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-installer/app/controller"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
)

// Services keys
const (
	UIControllerKey        = "service.controller.ui"
	InstallerControllerKey = "service.controller.installer"
)

func init() {
	service.Set(UIControllerKey, func(c service.Container) interface{} {
		return controller.NewUIController()
	})

	service.Set(InstallerControllerKey, func(c service.Container) interface{} {
		dockerClient := c.Get(DockerKey).(*docker.Client)

		return controller.NewInstallerController(dockerClient)
	})
}

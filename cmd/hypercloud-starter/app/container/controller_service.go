// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-starter/app/controller"
)

// Services keys
const (
	IndexControllerKey = "service.controller.index"
)

func init() {
	service.Set(IndexControllerKey, func(c service.Container) interface{} {
		return controller.NewIndexController()
	})
}

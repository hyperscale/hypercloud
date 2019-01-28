// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"github.com/euskadi31/go-server/request"
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-server/app/asset"
	"github.com/rs/zerolog/log"
)

// Services keys
const (
	ValidatorKey = "service.validator"
)

func init() {
	service.Set(ValidatorKey, func(c service.Container) interface{} {
		validator := request.NewValidator()

		if schema, err := asset.Asset("schema/service.json"); err == nil {
			if err := validator.AddSchemaFromJSON("service", schema); err != nil {
				log.Fatal().Err(err).Msg("alidator.AddSchemaFromJSON")
			}
		} else {
			log.Fatal().Err(err).Msg("Asset: schema/service.json")
		}

		if schema, err := asset.Asset("schema/stack.json"); err == nil {
			if err := validator.AddSchemaFromJSON("stack", schema); err != nil {
				log.Fatal().Err(err).Msg("alidator.AddSchemaFromJSON")
			}
		} else {
			log.Fatal().Err(err).Msg("Asset: schema/stack.json")
		}

		return validator
	})
}

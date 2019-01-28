// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"fmt"
	"strings"

	"github.com/asdine/storm"
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-server/app/config"
	"github.com/rs/zerolog/log"
)

// Services keys
const (
	DBKey = "service.database"
)

func init() {
	service.Set(DBKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)

		path := strings.TrimRight(cfg.Database.Path, "/")

		db, err := storm.Open(fmt.Sprintf("%s/hypercloud.db", path))
		if err != nil {
			log.Fatal().Err(err).Msg(DBKey)
		}

		return db // *storm.DB
	})
}

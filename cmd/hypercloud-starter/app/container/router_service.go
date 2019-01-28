// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"fmt"
	"net/http"

	server "github.com/euskadi31/go-server"
	"github.com/euskadi31/go-server/response"
	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-starter/app/config"
	hlogger "github.com/hyperscale/hypercloud/pkg/hypercloud/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// Services keys
const (
	RouterKey = "service.http.router"
)

func init() {
	service.Set(RouterKey, func(c service.Container) interface{} {
		cfg := c.Get(ConfigKey).(*config.Configuration)
		logger := c.Get(LoggerKey).(zerolog.Logger)
		indexController := c.Get(IndexControllerKey).(server.Controller)

		router := server.New(cfg.Server.ToConfig())

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(hlogger.Handler))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

		router.EnableCors()
		router.EnableHealthCheck()
		router.EnableRecovery()

		router.SetNotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
			response.Encode(w, r, http.StatusNotFound, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%s %s not found", r.Method, r.URL.Path),
				},
			})
		})

		router.AddController(indexController)

		return router // *server.Server
	})
}

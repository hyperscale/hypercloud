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
	"github.com/gorilla/handlers"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-server/app/config"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
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
		dockerClient := c.Get(DockerKey).(*docker.Client)
		topologyController := c.Get(TopologyControllerKey).(server.Controller)
		applicationController := c.Get(ServiceControllerKey).(server.Controller)
		deployController := c.Get(DeployControllerKey).(server.Controller)
		stackController := c.Get(StackControllerKey).(server.Controller)
		dockerController := c.Get(DockerControllerKey).(server.Controller)
		eventController := c.Get(EventControllerKey).(server.Controller)
		versionController := c.Get(VersionControllerKey).(server.Controller)

		router := server.New(cfg.Server.ToConfig())

		router.EnableCorsWithOptions(
			handlers.AllowedOrigins([]string{
				"*",
			}),
			handlers.AllowedHeaders([]string{
				"Authorization",
				"Content-Type",
				"X-Requested-With",
			}),
			handlers.AllowedMethods([]string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
				http.MethodOptions,
			}),
		)

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(hlogger.Handler))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

		router.EnableHealthCheck()
		router.EnableRecovery()

		router.SetNotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
			response.Encode(w, r, http.StatusNotFound, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%s %s not found", r.Method, r.URL.Path),
				},
			})
		})

		router.AddHealthCheck("docker", docker.NewHealthCheck(dockerClient))

		router.AddController(topologyController)
		router.AddController(applicationController)
		router.AddController(deployController)
		router.AddController(stackController)
		router.AddController(dockerController)
		router.AddController(eventController)
		router.AddController(versionController)

		return router // *server.Server
	})
}

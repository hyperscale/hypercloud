// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpaas/cmd/hyperpaas-installer/controller"
	"github.com/hyperscale/hyperpaas/config"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/version"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// const of service name
const (
	ServiceLoggerKey              string = "service.logger"
	ServiceConfigKey                     = "service.config"
	ServiceRouterKey                     = "service.router"
	ServiceDockerKey                     = "service.docker.client"
	ServiceAppKey                        = "service.app"
	ServiceInstallerControllerKey        = "service.controller.installer"
	ServiceUIControllerKey               = "service.controller.ui"
)

const applicationName = "hyperpaas-installer"

// Service Container
var container = service.New()

func init() {
	// Logger Service
	container.Set(ServiceLoggerKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		logger := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", cfg.Logger.Prefix).
			Str("version", version.Version.String()).
			Logger()

		zerolog.SetGlobalLevel(cfg.Logger.Level())

		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		stdlog.SetFlags(0)
		stdlog.SetOutput(logger)

		log.Logger = logger

		return logger
	})

	// Config Service
	container.Set(ServiceConfigKey, func(c *service.Container) interface{} {
		var cfgFile string
		cmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		cmd.StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")

		// Ignore errors; cmd is set for ExitOnError.
		cmd.Parse(os.Args[1:])

		options := viper.New()

		if cfgFile != "" { // enable ability to specify config file via flag
			options.SetConfigFile(cfgFile)
		}

		options.SetDefault("server.host", "")
		options.SetDefault("server.port", 8080)
		options.SetDefault("server.shutdown_timeout", 10*time.Second)
		options.SetDefault("server.write_timeout", 10*time.Second)
		options.SetDefault("server.read_timeout", 10*time.Second)
		options.SetDefault("server.read_header_timeout", 10*time.Millisecond)
		options.SetDefault("logger.level", "info")
		options.SetDefault("logger.prefix", applicationName)

		options.SetConfigName("config") // name of config file (without extension)

		options.AddConfigPath("/etc/" + applicationName + "/")   // path to look for the config file in
		options.AddConfigPath("$HOME/." + applicationName + "/") // call multiple times to add many search paths
		options.AddConfigPath(".")

		if port := os.Getenv("PORT"); port != "" {
			os.Setenv("HYPERPAAS_INSTALLER_PORT", port)
		}

		options.SetEnvPrefix("HYPERPAAS_INSTALLER")
		options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		options.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := options.ReadInConfig(); err == nil {
			log.Info().Msgf("Using config file: %s", options.ConfigFileUsed())
		}

		return config.NewConfiguration(options)
	})

	container.Set(ServiceDockerKey, func(c *service.Container) interface{} {
		dc, err := docker.NewEnvClient()
		if err != nil {
			log.Fatal().Err(err).Msg(ServiceDockerKey)
		}

		return dc
	})

	container.Set(ServiceUIControllerKey, func(c *service.Container) interface{} {
		controller, err := controller.NewUIController()
		if err != nil {
			log.Fatal().Err(err).Msg(ServiceUIControllerKey)
		}

		return controller
	})

	container.Set(ServiceInstallerControllerKey, func(c *service.Container) interface{} {
		dockerClient := c.Get(ServiceDockerKey).(*docker.Client)

		controller, err := controller.NewInstallerController(dockerClient)
		if err != nil {
			log.Fatal().Err(err).Msg(ServiceInstallerControllerKey)
		}

		return controller
	})

	// Router Service
	container.Set(ServiceRouterKey, func(c *service.Container) interface{} {
		logger := c.Get(ServiceLoggerKey).(zerolog.Logger)
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)
		dockerClient := c.Get(ServiceDockerKey).(*docker.Client)
		installerController := c.Get(ServiceInstallerControllerKey).(server.Controller)
		uiController := c.Get(ServiceUIControllerKey).(server.Controller)

		router := server.NewRouter()

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			rlog := hlog.FromRequest(r)

			var evt *zerolog.Event

			switch {
			case status >= 200 && status <= 299:
				evt = rlog.Info()
			case status >= 300 && status <= 399:
				evt = rlog.Info()
			case status >= 400 && status <= 499:
				evt = rlog.Warn()
			default:
				evt = rlog.Error()
			}

			evt.
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msgf("%s %s", r.Method, r.URL.Path)
		}))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

		router.EnableHealthCheck()
		router.EnableRecovery()
		router.EnableCorsWithOptions(cors.Options{
			AllowCredentials: true,
			AllowedOrigins:   []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodOptions,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
				"X-Requested-With",
			},
			Debug: cfg.Server.Debug,
		})

		router.SetNotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
			server.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%s %s not found", r.Method, r.URL.Path),
				},
			})
		})

		router.AddHealthCheck("docker", docker.NewHealthCheck(dockerClient))

		router.AddController(uiController)
		router.AddController(installerController)

		return router
	})
}

// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package container

import (
	"flag"
	"os"
	"strings"
	"time"

	service "github.com/euskadi31/go-service"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-server/app/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Services keys
const (
	ConfigKey = "service.config"
)

const name = "hypercloud-server"

func init() {
	service.Set(ConfigKey, func(c service.Container) interface{} {
		cmd := c.Get(FlagsKey).(*flag.FlagSet)

		var cfgFile string

		cmd.StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")

		// Ignore errors; cmd is set for ExitOnError.
		// nolint:gosec
		_ = cmd.Parse(os.Args[1:])

		cfg := config.NewConfiguration()

		options := viper.New()

		options.SetDefault("logger.level", "info")
		options.SetDefault("logger.prefix", name)
		options.SetDefault("server.http.host", "")
		options.SetDefault("server.http.port", 8080)
		options.SetDefault("server.profiling", true)
		options.SetDefault("server.metrics", true)
		options.SetDefault("server.healthcheck", true)
		options.SetDefault("server.shutdown_timeout", 10*time.Second)
		options.SetDefault("server.write_timeout", 0)
		options.SetDefault("server.read_timeout", 10*time.Second)
		options.SetDefault("server.read_header_timeout", 10*time.Millisecond)
		options.SetDefault("database.path", "/var/lib/hypercloud")
		options.SetDefault("docker.host", "unix:///var/run/docker.sock")

		options.SetConfigName("config") // name of config file (without extension)

		options.AddConfigPath("/etc/" + name + "/")   // path to look for the config file in
		options.AddConfigPath("$HOME/." + name + "/") // call multiple times to add many search paths
		options.AddConfigPath(".")

		if cfgFile != "" { // enable ability to specify config file via flag
			options.SetConfigFile(cfgFile)
		}

		if port := os.Getenv("PORT"); port != "" {
			os.Setenv("HYPERCLOUD_SERVER_SERVER_HTTP_PORT", port)
		}

		options.SetEnvPrefix("HYPERCLOUD_SERVER")
		options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		options.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := options.ReadInConfig(); err == nil {
			log.Info().Msgf("Using config file: %s", options.ConfigFileUsed())
		}

		if err := options.Unmarshal(cfg); err != nil {
			log.Fatal().Err(err).Msg(ConfigKey)
		}

		return cfg // *config.Configuration
	})
}

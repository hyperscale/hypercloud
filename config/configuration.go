// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/spf13/viper"
)

// Configuration struct
type Configuration struct {
	Logger   *LoggerConfiguration
	Server   *ServerConfiguration
	Docker   *DockerConfiguration
	Database *DatabaseConfiguration
}

// NewConfiguration constructor
func NewConfiguration(options *viper.Viper) *Configuration {
	return &Configuration{
		Server: &ServerConfiguration{
			Host:              options.GetString("server.host"),
			Port:              options.GetInt("server.port"),
			Debug:             options.GetBool("server.debug"),
			ShutdownTimeout:   options.GetDuration("server.shutdown_timeout"),
			WriteTimeout:      options.GetDuration("server.write_timeout"),
			ReadTimeout:       options.GetDuration("server.read_timeout"),
			ReadHeaderTimeout: options.GetDuration("server.read_header_timeout"),
		},
		Logger: &LoggerConfiguration{
			LevelName: options.GetString("logger.level"),
			Prefix:    options.GetString("logger.prefix"),
		},
		Docker: &DockerConfiguration{
			Host: options.GetString("docker.host"),
		},
		Database: &DatabaseConfiguration{
			Path: options.GetString("database.path"),
		},
	}
}

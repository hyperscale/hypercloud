// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

// Configuration struct
type Configuration struct {
	Debug      bool   `env:"DEBUG" envDefault:"false"`
	Port       int    `env:"PORT" envDefault:"8080"`
	DockerHost string `env:"DOCKER_HOST" envDefault:"unix:///var/run/docker.sock"`
	Path       string `env:"HYPERPAAS_PATH" envDefault:"/var/lib/hyperpaas"`
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/hyperscale/hypercloud/pkg/hypercloud/logger"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/server"
)

// Configuration struct
type Configuration struct {
	Logger *logger.Configuration
	Server *server.Configuration
}

// NewConfiguration constructor
func NewConfiguration() *Configuration {
	return &Configuration{
		Server: &server.Configuration{},
		Logger: &logger.Configuration{},
	}
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docker

import (
	"context"

	server "github.com/euskadi31/go-server"
)

// NewHealthCheck handle
func NewHealthCheck(d *Client) server.HealthCheckHandlerFunc {
	return server.HealthCheckHandlerFunc(func() bool {
		if _, err := d.Ping(context.Background()); err != nil {
			return false
		}

		return true
	})
}

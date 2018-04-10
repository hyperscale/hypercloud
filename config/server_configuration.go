// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"net"
	"strconv"
	"time"
)

// ServerConfiguration struct
type ServerConfiguration struct {
	Host              string
	Port              int
	Debug             bool
	ShutdownTimeout   time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

// Addr string
func (c ServerConfiguration) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http/httputil"

	server "github.com/euskadi31/go-server"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/docker"
)

// DockerController struct
type DockerController struct {
	proxy *httputil.ReverseProxy
}

// NewDockerController func
func NewDockerController(host string) (*DockerController, error) {
	proxy, err := docker.NewProxy(host, "/docker/")
	if err != nil {
		return nil, err
	}

	return &DockerController{
		proxy: proxy,
	}, nil
}

// Mount endpoints
func (c DockerController) Mount(r *server.Router) {
	r.PathPrefix("/docker/").Handler(c.proxy)
}

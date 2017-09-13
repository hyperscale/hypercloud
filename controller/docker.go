// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http/httputil"

	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/server"
)

type DockerController struct {
	proxy *httputil.ReverseProxy
}

func NewDockerController(host string) (*DockerController, error) {
	proxy, err := docker.NewProxy(host)
	if err != nil {
		return nil, err
	}

	return &DockerController{
		proxy: proxy,
	}, nil
}

func (c DockerController) Mount(r *server.Router) {
	r.AddPrefixRoute("/api/", c.proxy)
}

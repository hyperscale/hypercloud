// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	server "github.com/euskadi31/go-server"
	"github.com/hyperscale/hypercloud/cmd/hypercloud-starter/app/asset"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/memfs"
)

type indexController struct {
}

// NewIndexController func
func NewIndexController() server.Controller {
	return &indexController{}
}

// Mount endpoints
func (c indexController) Mount(r *server.Router) {
	r.HandleFunc("/", c.getIndexHandler).Methods(http.MethodGet)
}

// GET /
func (c indexController) getIndexHandler(w http.ResponseWriter, r *http.Request) {
	name := "static/index.html"

	w.Header().Set("Retry-After", "3600")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)

	body, _ := asset.Asset(name)

	info, _ := asset.AssetInfo(name)

	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		memfs.NewBuffer(&body),
	)
}

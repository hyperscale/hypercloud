// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpaas/cmd/hyperpaas-installer/assets"
	"github.com/hyperscale/hyperpaas/memfs"
)

// UIController struct
type UIController struct {
}

// NewUIController func
func NewUIController() (*UIController, error) {
	return &UIController{}, nil
}

// Mount endpoints
func (c UIController) Mount(r *server.Router) {
	r.AddRouteFunc("/", c.getUIHandler).Methods(http.MethodGet)
}

func (c UIController) getUIHandler(w http.ResponseWriter, r *http.Request) {
	name := "static/index.html"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	body, err := assets.Asset(name)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	info, err := assets.AssetInfo(name)
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	http.ServeContent(
		w,
		r,
		info.Name(),
		info.ModTime(),
		memfs.NewBuffer(&body),
	)
}

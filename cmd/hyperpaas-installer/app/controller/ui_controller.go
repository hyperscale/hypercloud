// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/euskadi31/go-server"
	"github.com/hyperscale/hyperpaas/cmd/hyperpaas-installer/app/asset"
	"github.com/hyperscale/hyperpaas/pkg/hyperpaas/memfs"
)

type uiController struct {
}

// NewUIController func
func NewUIController() server.Controller {
	return &uiController{}
}

// Mount endpoints
func (c uiController) Mount(r *server.Router) {
	r.HandleFunc("/", c.getUIHandler).Methods(http.MethodGet)
}

func (c uiController) getUIHandler(w http.ResponseWriter, r *http.Request) {
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

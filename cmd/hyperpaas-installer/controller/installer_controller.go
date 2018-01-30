// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-sse"
	"github.com/hyperscale/hyperpaas/cmd/hyperpaas-installer/assets"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/docker/compose"
	"github.com/rs/zerolog/log"
)

// Event struct
type Event struct {
	Type string `json:"type"`
}

// InstallerController struct
type InstallerController struct {
	dockerClient *docker.Client
	events       chan interface{}
}

// NewInstallerController func
func NewInstallerController(dockerClient *docker.Client) (*InstallerController, error) {
	return &InstallerController{
		dockerClient: dockerClient,
		events:       make(chan interface{}, 10),
	}, nil
}

// Mount endpoints
func (c InstallerController) Mount(r *server.Router) {
	events := sse.NewServer(c.getEventsHandler)
	events.SetRetry(time.Second * 5)

	r.AddRoute("/installer/events", events).Methods(http.MethodGet)
	r.AddRouteFunc("/installer", c.postInstallerHandler).Methods(http.MethodPost)
}

func (c *InstallerController) postInstallerHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2048); err != nil {
		server.FailureFromError(w, http.StatusRequestEntityTooLarge, err)

		return
	}

	log.Info().Msgf("Data: %#v", r.PostForm)

	c.events <- Event{
		Type: "started",
	}

	content, err := assets.Asset("static/config/docker-compose.yml")
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	composeConfig, err := compose.Loader(content, map[string]string{
		"EMAIL":  r.PostForm.Get("email"),
		"DOMAIN": r.PostForm.Get("domain"),
	})
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	c.events <- Event{
		Type: "configured",
	}

	//@TODO exec deploy stack

	c.events <- Event{
		Type: "finished",
	}

	server.JSON(w, http.StatusOK, composeConfig)
}

func (c *InstallerController) getEventsHandler(rw sse.ResponseWriter, r *http.Request) {
	for {
		select {
		case event := <-c.events:
			data, err := json.Marshal(event)
			if err != nil {
				log.Error().Err(err).Msg("Marshal Event")

				continue
			}

			rw.Send(&sse.MessageEvent{
				Data: data,
			})

		case <-rw.CloseNotify:

			return
		}
	}
}

// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coreos/go-semver/semver"

	server "github.com/euskadi31/go-server"
	"github.com/euskadi31/go-server/response"
	sse "github.com/euskadi31/go-sse"
	"github.com/hyperscale/hypercloud/pkg/hypercloud/version"
	"github.com/rs/zerolog/log"
)

// VersionResponse struct
type VersionResponse struct {
	Current *semver.Version `json:"current"`
	Latest  *semver.Version `json:"latest"`
}

// VersionController struct
type VersionController struct {
	latest *semver.Version
}

// NewVersionController func
func NewVersionController() (*VersionController, error) {
	c := &VersionController{
		latest: version.Version,
	}

	go c.worker()

	return c, nil
}

// Mount endpoints
func (c VersionController) Mount(r *server.Router) {
	events := sse.NewServer(c.getVersionEventHandler)
	events.SetRetry(time.Minute * 30)

	r.Handle("/v1/version/latest", events).Methods(http.MethodGet).Headers("Accept", "text/event-stream")
	r.HandleFunc("/v1/version/latest", c.getVersionHandler).Methods(http.MethodGet)
}

func (c VersionController) fetcher() {
	v, err := version.GetLatest()
	if err != nil {
		log.Error().Err(err).Msg("Get latest version")

		return
	}

	c.latest = v

	if version.Version.LessThan(*v) {
		log.Info().Msgf("New Version Available: %s", v.String())
	}
}

func (c *VersionController) worker() {
	tickChan := time.NewTicker(time.Hour * 1).C

	go c.fetcher()

	for range tickChan {
		c.fetcher()
	}
}

// swagger:route GET /v1/version/latest Version getVersionHandler
//
// Get the latest version
//
//     Responses:
//       200: Version
//
func (c VersionController) getVersionHandler(w http.ResponseWriter, r *http.Request) {
	response.Encode(w, r, http.StatusOK, VersionResponse{
		Current: version.Version,
		Latest:  c.latest,
	})
}

// swagger:route GET /v1/version/latest Version getVersionEventHandler
//
// Get the latest version
//
//     Responses:
//       200: Version
//
func (c VersionController) getVersionEventHandler(rw sse.ResponseWriter, r *http.Request) {
	latestNotified := version.Version

	tickChan := time.NewTicker(time.Minute * 1).C

	for {
		select {
		case <-tickChan:
			if !latestNotified.Equal(*c.latest) {
				data, err := json.Marshal(VersionResponse{
					Latest: c.latest,
				})
				if err != nil {
					log.Error().Err(err).Msg("Marshal VersionResponse")

					continue
				}

				rw.Send(&sse.MessageEvent{
					Data: data,
				})

				latestNotified = c.latest
			}

		case <-r.Context().Done():
			return
		}
	}
}

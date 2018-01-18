// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/euskadi31/go-server"
	sse "github.com/euskadi31/go-sse"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/rs/zerolog/log"
)

// EventController struct
type EventController struct {
	dockerClient *docker.Client
}

// NewEventController func
func NewEventController(dockerClient *docker.Client) (*EventController, error) {
	return &EventController{
		dockerClient: dockerClient,
	}, nil
}

// Mount endpoints
func (c EventController) Mount(r *server.Router) {
	events := sse.NewServer(c.getEventsHandler)
	events.SetRetry(time.Second * 5)

	r.AddRoute("/v1/events", events).Methods(http.MethodGet)
}

// swagger:route GET /v1/events Event getEventsHandler
//
// Get the event list
//
//     Responses:
//       200: Event
//
func (c EventController) getEventsHandler(rw sse.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// recovery
	lastID := r.Header.Get("Last-Event-ID")
	if lastID != "" {
		log.Printf("Recovery with ID: %s\n", lastID)
	}

	eventsCh, errCh := c.dockerClient.Events(ctx, types.EventsOptions{})

	for {
		select {
		case event := <-eventsCh:
			data, err := json.Marshal(event)
			if err != nil {
				log.Error().Err(err).Msg("Marshal Event")

				continue
			}

			rw.Send(&sse.MessageEvent{
				ID:    strconv.FormatInt(event.TimeNano, 10),
				Event: event.Type,
				Data:  data,
			})
		case err := <-errCh:
			log.Error().Err(err).Msg("Events")

		case <-rw.CloseNotify:

			return
		}
	}
}

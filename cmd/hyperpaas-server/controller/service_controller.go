// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-sse"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/hyperscale/hyperpaas/database/entity"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/http/request"
	"github.com/hyperscale/hyperpaas/http/response"
	"github.com/rs/zerolog/log"
)

type ServiceHosts map[string]bool

func (h ServiceHosts) Add(host string) {
	if _, ok := h[host]; ok {
		return
	}

	h[host] = true
}

func (h ServiceHosts) String() string {
	hosts := []string{}

	for host := range h {
		hosts = append(hosts, host)
	}

	return fmt.Sprintf("Host:%s", strings.Join(hosts, ","))
}

// ServiceController struct
type ServiceController struct {
	dockerClient *docker.Client
	db           *storm.DB
	validator    *server.Validator
	queryDecoder *schema.Decoder
}

// NewServiceController func
func NewServiceController(dockerClient *docker.Client, db *storm.DB, validator *server.Validator) (*ServiceController, error) {
	if err := db.Init(&entity.Service{}); err != nil {
		return nil, err
	}

	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	return &ServiceController{
		dockerClient: dockerClient,
		db:           db,
		validator:    validator,
		queryDecoder: decoder,
	}, nil
}

// Mount endpoints
func (c ServiceController) Mount(r *server.Router) {
	events := sse.NewServer(c.getServiceStatsHandler)
	events.SetRetry(time.Second * 5)

	r.AddRouteFunc("/v1/services", c.getServicesHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/v1/services", c.postServiceHandler).Methods(http.MethodPost)
	r.AddRouteFunc("/v1/services/{id:[0-9a-z]{25}}", c.getServiceHandler).Methods(http.MethodGet)
	r.AddRouteFunc("/v1/services/{id:[0-9a-z]{25}}", c.putServiceHandler).Methods(http.MethodPut)
	r.AddRoute("/v1/services/{id:[0-9a-z]{25}}/stats", events).Methods(http.MethodGet)
}

// swagger:route GET /v1/services Service getServicesHandler
//
// Get the services list
//
//     Responses:
//       200: Service
//
func (c ServiceController) getServicesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := &request.ServicesQuery{}

	if err := c.queryDecoder.Decode(query, r.URL.Query()); err != nil {
		log.Error().Err(err).Msg("Decode query parameters")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	filter := filters.NewArgs()

	if query.StackID != "" {
		filter.Add("label", "com.docker.stack.namespace="+query.StackID)
	}

	services, err := c.dockerClient.ServiceList(ctx, types.ServiceListOptions{
		Filters: filter,
	})
	if err != nil {
		log.Error().Err(err).Msg("ServiceList")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	response := []swarm.Service{}

	for _, service := range services {
		labels := service.Spec.Labels

		if _, ok := labels[docker.LabelInternalNamespace]; ok {
			continue
		}

		response = append(response, service)
	}

	server.JSON(w, http.StatusOK, response)
}

// swagger:route POST /v1/services Service postServiceHandler
//
// Create service
//     Request: ServiceCreateRequest
//     Responses:
//       201: Service
//
func (c ServiceController) postServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	service := &request.ServiceCreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(service); err != nil {
		log.Error().Err(err).Msg("Decode body request")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	if result := c.validator.Validate("service", service); !result.IsValid() {
		log.Error().Err(result.AsError()).Msg("Validate body request")

		server.FailureFromValidator(w, result)

		return
	}

	hosts := ServiceHosts{}

	hosts.Add(fmt.Sprintf("%s.%s", service.Name, "hyperpaas.service"))

	spec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: service.GetServiceName(),
			Labels: map[string]string{
				docker.LabelStackNamespace:     service.StackID,
				"traefik.docker.network":       "traefik-net",
				"traefik.enable":               "true",
				"traefik.frontend.entryPoints": "https",
				"traefik.frontend.rule":        hosts.String(),
				"traefik.port":                 "80",
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: "dockercloud/hello-world:latest",
				Labels: map[string]string{
					docker.LabelStackNamespace: service.StackID,
				},
				Env: []string{
					"PORT=80",
				},
				// Isolation: container.IsolationDefault,
			},
		},
	}

	resp, err := c.dockerClient.ServiceCreate(ctx, spec, types.ServiceCreateOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Docker Service Create")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	log.Info().Msgf("Service ID: %s", resp.ID)

	for _, w := range resp.Warnings {
		log.Warn().Msg(w)
	}

	// service.Hosts = append(service.Hosts, fmt.Sprintf("%s.%s", service.Name, "hyperpaas.service"))

	server.JSON(w, http.StatusCreated, response.ServiceCreateResponse{
		ServiceCreateRequest: service,
		ID:                   resp.ID,
	})
}

// swagger:route GET /v1/services/{id} Stack getServiceHandler
//
// Get a service by id
//
//     Responses:
//       200: Service
//
func (c ServiceController) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)

	id := params["id"]

	log.Debug().Msgf("Service ID: %s", id)

	service, _, err := c.dockerClient.ServiceInspectWithRaw(ctx, id, types.ServiceInspectOptions{})
	if err != nil {
		server.FailureFromError(w, http.StatusNotFound, err)

		return
	}

	server.JSON(w, http.StatusOK, service)
}

// swagger:route GET /v1/services/{id}/stats Stack getServiceHandler
//
// Get stats for service by id
//
//     Responses:
//       200: Service
//
func (c ServiceController) getServiceStatsHandler(rw sse.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)

	id := params["id"]

	log.Debug().Msgf("Service ID: %s", id)

	service, _, err := c.dockerClient.ServiceInspectWithRaw(ctx, id, types.ServiceInspectOptions{})
	if err != nil {
		log.Error().Err(err).Msg("ServiceInspectWithRaw")

		return
	}

	filter := filters.NewArgs()
	filter.Add("service", service.Spec.Name)
	filter.Add("desired-state", "running")

	tasks, err := c.dockerClient.TaskList(ctx, types.TaskListOptions{
		Filters: filter,
	})
	if err != nil {
		log.Error().Err(err).Msg("TaskList")

		return
	}

	containers := []string{}

	/*
		swarm.Task{
			ID:"mosbb0aw0oobw9atkrh2pjooj",
			Meta:swarm.Meta{
				Version:swarm.Version{Index:0xd2},
				CreatedAt:time.Time{wall:0x6c64f62, ext:63651657051, loc:(*time.Location)(nil)},
				UpdatedAt:time.Time{wall:0x1a3a8ec6, ext:63651657059, loc:(*time.Location)(nil)}
			},
			Annotations:swarm.Annotations{Name:"", Labels:map[string]string{}},
			Spec:swarm.TaskSpec{
				ContainerSpec:(*swarm.ContainerSpec)(0xc420270160),
				PluginSpec:(*runtime.PluginSpec)(nil),
				Resources:(*swarm.ResourceRequirements)(0xc4201cc060),
				RestartPolicy:(*swarm.RestartPolicy)(0xc4203f6720),
				Placement:(*swarm.Placement)(0xc42010c050),
				Networks:[]swarm.NetworkAttachmentConfig{
					swarm.NetworkAttachmentConfig{
						Target:"lpq4pehw4foqeae2nlm2n6r05",
						Aliases:[]string{"web"},
						DriverOpts:map[string]string(nil)
					}
				},
				LogDriver:(*swarm.Driver)(nil),
				ForceUpdate:0x0,
				Runtime:""
			},
			ServiceID:"pll7obzsx6aksevjl0oery6dq",
			Slot:1,
			NodeID:"g3z3mv1eunnqbgn1p1t9elfoq",
			Status:swarm.TaskStatus{Timestamp:time.Time{wall:0x15c214f0, ext:63651657059, loc:(*time.Location)(nil)},
			State:"running",
			Message:"started",
			Err:"",
			ContainerStatus:swarm.ContainerStatus{
				ContainerID:"54cc8bed986caf955c976006d603fe7e922c154cb5509b8c64797ee135aa3aea",
				PID:2777,
				ExitCode:0
			},
			PortStatus:swarm.PortStatus{Ports:[]swarm.PortConfig(nil)}}, DesiredState:"running", NetworksAttachments:[]swarm.NetworkAttachment{swarm.NetworkAttachment{Network:swarm.Network{ID:"lpq4pehw4foqeae2nlm2n6r05", Meta:swarm.Meta{Version:swarm.Version{Index:0xbc}, CreatedAt:time.Time{wall:0x17312498, ext:63651391082, loc:(*time.Location)(nil)}, UpdatedAt:time.Time{wall:0x38b511d, ext:63651657049, loc:(*time.Location)(nil)}},Spec:swarm.NetworkSpec{Annotations:swarm.Annotations{Name:"acme_frontend", Labels:map[string]string{"com.docker.stack.namespace":"acme"}}, DriverConfiguration:(*swarm.Driver)(0xc4200e2700), IPv6Enabled:false, Internal:false, Attachable:false, Ingress:false, IPAMOptions:(*swarm.IPAMOptions)(nil), ConfigFrom:(*network.ConfigReference)(nil), Scope:"swarm"}, DriverState:swarm.Driver{Name:"overlay", Options:map[string]string{"com.docker.network.driver.overlay.vxlanid_list":"4097"}}, IPAMOptions:(*swarm.IPAMOptions)(0xc4203f7050)}, Addresses:[]string{"10.0.0.6/24"}}}, GenericResources:[]swarm.GenericResource(nil)}
	*/

	for _, task := range tasks {
		// log.Debug().Msgf("Task: %#v", task.ID)

		containers = append(containers, task.Status.ContainerStatus.ContainerID)
	}

	events := make(chan json.RawMessage, len(containers)*2)

	for _, container := range containers {
		go func(container string) {
			containerStats, err := c.dockerClient.ContainerStats(ctx, container, true)
			if err != nil {
				log.Error().Err(err).Msg("ContainerStats")

				return
			}

			defer containerStats.Body.Close()

			decoder := json.NewDecoder(containerStats.Body)

			for {
				//var stats types.StatsJSON
				var stats json.RawMessage

				if err := decoder.Decode(&stats); err == io.EOF {
					break
				} else if err != nil {
					log.Error().Err(err).Msg("JSON Decoder")

					return
				}

				events <- stats
			}
		}(container)
	}

	for {
		select {
		case stats := <-events:
			rw.Send(&sse.MessageEvent{
				Data: stats,
			})
		case <-rw.CloseNotify:
			close(events)

			return
		}
	}
}

// swagger:route POST /v1/services/{id} Service putServiceHandler
//
// Update service
//     Request: Service
//     Responses:
//       200: Service
//
func (c ServiceController) putServiceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)

	id := params["id"]

	log.Debug().Msgf("Service ID: %s", id)

	var service swarm.Service

	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		log.Error().Err(err).Msg("Unmarshal Request Body")

		server.FailureFromError(w, http.StatusBadRequest, err)

		return
	}

	resp, err := c.dockerClient.ServiceUpdate(ctx, id, service.Version, service.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Docker Service Update")

		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	for _, w := range resp.Warnings {
		log.Warn().Msg(w)
	}

	server.JSON(w, http.StatusOK, service)
}

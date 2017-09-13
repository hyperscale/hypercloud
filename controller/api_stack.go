// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"net/http"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/hyperscale/hyperpaas/docker"
	"github.com/hyperscale/hyperpaas/entity"
	"github.com/hyperscale/hyperpaas/server"
	"github.com/pkg/errors"
)

const (
	// LabelNamespace is the label used to track stack resources
	LabelNamespace = "com.docker.stack.namespace"
)

type byName []*entity.Stack

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }

// GetStacksHandler /api/stacks
func (c APIController) GetStacksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	dc, err := docker.FromContext(r.Context())
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	services, err := dc.ServiceList(ctx, types.ServiceListOptions{Filters: getAllStacksFilter()})
	if err != nil {
		server.FailureFromError(w, http.StatusInternalServerError, err)

		return
	}

	m := make(map[string]*entity.Stack, 0)
	for _, service := range services {
		labels := service.Spec.Labels
		name, ok := labels[LabelNamespace]
		if !ok {
			server.FailureFromError(w, http.StatusInternalServerError, errors.Errorf("cannot get label %s for service %s", LabelNamespace, service.ID))

			return
		}

		ztack, ok := m[name]
		if !ok {
			m[name] = &entity.Stack{
				Name:     name,
				Services: 1,
			}
		} else {
			ztack.Services++
		}
	}

	stacks := []*entity.Stack{}

	for _, stack := range m {
		stacks = append(stacks, stack)
	}

	sort.Sort(byName(stacks))

	server.JSON(w, http.StatusOK, stacks)
}

func getAllStacksFilter() filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", LabelNamespace)

	return filter
}

// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docker

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"
)

const (
	// LabelStackNamespace is the label used to track stack resources
	LabelStackNamespace = "com.docker.stack.namespace"
)

type byName []Stack

func (n byName) Len() int           { return len(n) }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }

// StackList returns the list of stacks.
func (c *Client) StackList(ctx context.Context) ([]Stack, error) {
	filter := filters.NewArgs()
	filter.Add("label", LabelStackNamespace)

	services, err := c.ServiceList(ctx, types.ServiceListOptions{
		Filters: filter,
	})
	if err != nil {
		return nil, err
	}

	m := make(map[string]Stack, 0)
	for _, service := range services {
		labels := service.Spec.Labels
		name, ok := labels[LabelStackNamespace]
		if !ok {
			return nil, errors.Errorf("cannot get label %s for service %s", LabelStackNamespace, service.ID)
		}

		ztack, ok := m[name]
		if !ok {
			m[name] = Stack{
				Name:     name,
				Services: 1,
			}
		} else {
			ztack.Services++
		}
	}

	stacks := []Stack{}

	for _, stack := range m {
		stacks = append(stacks, stack)
	}

	sort.Sort(byName(stacks))

	return stacks, nil
}

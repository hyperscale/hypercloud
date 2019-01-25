// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	std "github.com/euskadi31/go-std"
)

// Service represents a Service.
//
// swagger:model Service
type Service struct {

	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
	ID string `storm:"id" json:"id"`

	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
	StackID string `storm:"index" json:"stack_id"`

	// Service name.
	//
	// required: true
	// min length: 3
	// pattern: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	Name string `json:"name" storm:"index"`

	CreatedAt std.DateTime `json:"created_at"`

	UpdatedAt std.DateTime `json:"updated_at,omitempty"`

	Hosts []string `json:"hosts"`
}

// swagger:response Service
//nolint
type serviceResponseDoc struct {
	// in: body
	Body Service
}

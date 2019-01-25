// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	std "github.com/euskadi31/go-std"
)

// StackCompose represents a Stack Compose.
//
// swagger:model StackCompose
type StackCompose struct {

	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
	ID string `storm:"id" json:"id"`

	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
	StackID string `storm:"index" json:"stack_id"`

	ComposeFile string `json:"-"`

	Version int `storm:"index" json:"version"`

	CreatedAt std.DateTime `json:"created_at"`

	UpdatedAt std.DateTime `json:"updated_at,omitempty"`
}

// swagger:response StackCompose
//nolint
type stackComposeResponseDoc struct {
	// in: body
	Body StackCompose
}

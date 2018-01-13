// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	std "github.com/euskadi31/go-std"
)

// Stack represents a Stack.
//
// swagger:model Stack
type Stack struct {

	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
	ID string `storm:"id" json:"id"`

	// Stack name.
	//
	// required: true
	// min length: 3
	// pattern: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	Name string `json:"name" storm:"index"`

	Services int `json:"services"`

	Version int `json:"version"`

	CreatedAt std.DateTime `json:"created_at"`
}

// swagger:response Stack
type stackResponseDoc struct {
	// in: body
	Body Stack
}

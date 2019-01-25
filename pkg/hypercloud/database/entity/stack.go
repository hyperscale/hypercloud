// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

// Stack represents a Stack.
//
// swagger:model Stack
type Stack struct {
	// Stack name.
	//
	// required: true
	// min length: 3
	// pattern: ^[a-z0-9]+(?:-[a-z0-9]+)*$
	Name string `json:"Name" storm:"id"`

	Services int `json:"Services"`
}

// swagger:response Stack
//nolint
type stackResponseDoc struct {
	// in: body
	Body Stack
}

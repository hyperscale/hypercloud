// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"regexp"
	"strings"

	std "github.com/euskadi31/go-std"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Application represents a Application.
//
// swagger:model Application
type Application struct {

	// pattern: ^[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89ABab][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12}$
	ID string `storm:"id" json:"id"`

	// Application name.
	//
	// required: true
	// min length: 3
	// pattern: ^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$
	Name string `json:"name"`

	CreatedAt std.DateTime `json:"created_at"`

	UpdatedAt std.DateTime `json:"updated_at,omitempty"`

	Hosts []string `json:"hosts"`
}

// Slug name
func (a Application) Slug() string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(a.Name), "-"), "-")
}

// swagger:response Application
type applicationResponseDoc struct {
	// in: body
	Body Application
}

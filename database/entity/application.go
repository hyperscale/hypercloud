// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

import (
	"regexp"
	"strings"

	"github.com/euskadi31/go-std"
)

var re = regexp.MustCompile("[^a-z0-9]+")

// Application struct
type Application struct {
	ID        string       `storm:"id" json:"id"`
	Name      string       `json:"name"`
	CreatedAt std.DateTime `json:"created_at"`
	UpdatedAt std.DateTime `json:"updated_at,omitempty"`
	Hosts     []string     `json:"hosts"`
}

// Slug name
func (a Application) Slug() string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(a.Name), "-"), "-")
}

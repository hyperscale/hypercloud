// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

// Registry entity
type Registry struct {
	ID       int    `json:"ID" storm:"id,increment"`
	Server   string `json:"Server" storm:"unique"`
	Username string `json:"Username"`
	Password string `json:"-"`
	Token    string `json:"Token"`
}

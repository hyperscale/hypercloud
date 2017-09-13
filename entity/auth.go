// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

// Auth entity
type Auth struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

// Validate auth entity
func (a Auth) Validate() bool {
	return a.Email != "" && a.Password != ""
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package entity

// User entity
type User struct {
	ID        int    `storm:"id,increment"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email" storm:"unique"`
	Password  string `json:"password"`
	Salt      string `json:"-"`
}

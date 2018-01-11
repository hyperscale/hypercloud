// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

var dependencies = map[string]struct {
	Image string
	Tag   string
}{
	"traefik": {
		Image: "traefik",
		Tag:   "1.4.6",
	},
	"registry": {
		Image: "registry",
		Tag:   "2.6.2",
	},
}

// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docker

import (
	"context"
)

// NewHealthCheck handle
func NewHealthCheck() func(ctx context.Context) bool {
	return func(ctx context.Context) bool {
		d, err := FromContext(ctx)
		if err != nil {
			return false
		}

		if _, err := d.Ping(context.Background()); err != nil {
			return false
		}

		return true
	}
}

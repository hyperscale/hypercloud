// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package request

// ServicesRequest struct
// swagger:parameters getServicesHandler
type ServicesRequest struct {
	// pattern: ^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$
	StackID string `schema:"stack_id"`
}

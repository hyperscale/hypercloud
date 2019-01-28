// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package request

// ServiceCreateRequest struct
// swagger:model
type ServiceCreateRequest struct {
	// pattern: ^[a-z][a-z0-9_]+$
	StackID string `json:"stack_id"`

	// pattern: ^[a-z][a-z0-9_]+$
	Name string `json:"name"`
}

// GetServiceName with stack_id
func (r ServiceCreateRequest) GetServiceName() string {
	return r.StackID + "_" + r.Name
}

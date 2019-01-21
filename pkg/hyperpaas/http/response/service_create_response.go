// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package response

import "github.com/hyperscale/hyperpaas/http/request"

// ServiceCreateResponse struct
type ServiceCreateResponse struct {
	*request.ServiceCreateRequest
	ID string `json:"id"`
}

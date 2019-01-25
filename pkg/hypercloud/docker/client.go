// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docker

import "github.com/docker/docker/client"

// Client struct
type Client struct {
	*client.Client
}

// NewEnvClient initializes a new API client based on environment variables.
func NewEnvClient() (*Client, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv)
	//dc, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: dc,
	}, nil
}

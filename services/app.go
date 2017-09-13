// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import "github.com/rs/zerolog/log"

// AppService struct
type AppService struct {
}

// NewAppService constructor
func NewAppService() *AppService {
	return &AppService{}
}

// Run App
func (s *AppService) Run() error {
	log.Info().Msg("hello")

	return nil
}

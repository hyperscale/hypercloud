// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import "github.com/rs/zerolog"

// AppService struct
type AppService struct {
	log zerolog.Logger
}

// NewAppService constructor
func NewAppService(log zerolog.Logger) *AppService {
	return &AppService{
		log: log,
	}
}

// Run App
func (s *AppService) Run() error {

	s.log.Info().Msg("hello")

	return nil
}

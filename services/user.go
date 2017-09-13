// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import (
	"github.com/asdine/storm"
	"github.com/hyperscale/hyperpaas/entity"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserService struct
type UserService struct {
	db *storm.DB
}

// NewUserService func
func NewUserService(db *storm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// Create user
func (s UserService) Create(user *entity.User) (*entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hash)

	if err := s.db.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID func
func (s UserService) GetByID(id int) (*entity.User, error) {
	user := &entity.User{}

	if err := s.db.One("ID", id, user); err != nil {
		return nil, errors.Wrapf(err, "Cannot find user by ID: %d", id)
	}

	return user, nil
}

// FindByEmail func
func (s UserService) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}

	if err := s.db.One("Email", email, user); err != nil {
		return nil, errors.Wrapf(err, "Cannot find user by email: %s", email)
	}

	return user, nil
}

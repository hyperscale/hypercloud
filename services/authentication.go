// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/hyperscale/hyperpaas/entity"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	privateKey = "private.pem"
	publicKey  = "public.pem"
)

// AuthenticationService struct
type AuthenticationService struct {
	path        string
	userService *UserService
	privateKey  []byte
	publicKey   []byte
}

// NewAuthenticationService func
func NewAuthenticationService(path string, userService *UserService) *AuthenticationService {
	return &AuthenticationService{
		path:        strings.TrimRight(path, "/"),
		userService: userService,
	}
}

func (s AuthenticationService) savePrivatePEMKey(filename string, key *rsa.PrivateKey) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", s.path, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	var privateKey = &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	}

	return pem.Encode(file, privateKey)
}

func (s AuthenticationService) savePublicPEMKey(filename string, pubkey rsa.PublicKey) error {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	//asn1Bytes, err := asn1.Marshal(pubkey)
	if err != nil {
		return err
	}

	var pemkey = &pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   asn1Bytes,
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", s.path, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, pemkey)
}

// HasKey exists
func (s *AuthenticationService) HasKey() bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", s.path, privateKey)); os.IsNotExist(err) {
		return false
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s", s.path, publicKey)); os.IsNotExist(err) {
		return false
	}

	return true
}

// LoadKeys files
func (s *AuthenticationService) LoadKeys() error {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", s.path, privateKey))
	if err != nil {
		return err
	}

	s.privateKey = bytes

	bytes, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", s.path, publicKey))
	if err != nil {
		return err
	}

	s.publicKey = bytes

	return nil
}

// GenerateKey files in path
func (s *AuthenticationService) GenerateKey() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	if err := s.savePrivatePEMKey(privateKey, key); err != nil {
		return err
	}

	if err := s.savePublicPEMKey(publicKey, key.PublicKey); err != nil {
		return err
	}

	return nil
}

// GenerateToken string
func (s *AuthenticationService) GenerateToken(claims jws.Claims) (string, error) {
	rsaPrivate, err := crypto.ParseRSAPrivateKeyFromPEM(s.privateKey)
	if err != nil {
		return "", err
	}

	//jwt := jws.NewJWT(claims, crypto.SigningMethodHS256)
	jwt := jws.NewJWT(claims, crypto.SigningMethodRS256)

	b, err := jwt.Serialize(rsaPrivate)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Validate token
func (s *AuthenticationService) Validate(token string) bool {
	rsaPublic, err := crypto.ParseRSAPublicKeyFromPEM(s.publicKey)
	if err != nil {
		log.Error().Err(err).Msg("")

		return false
	}

	jwt, err := jws.ParseJWT([]byte(token))
	if err != nil {
		log.Error().Err(err).Msg("")

		return false
	}

	// Validate token
	if err = jwt.Validate(rsaPublic, crypto.SigningMethodRS256); err != nil {
		log.Error().Err(err).Msg("")

		return false
	}

	return true
}

// Authenticate user
func (s *AuthenticationService) Authenticate(req *entity.Auth) (*entity.Token, error) {
	if !req.Validate() {
		return nil, errors.New("username or password is incorrect")
	}

	user, err := s.userService.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	expiration := time.Now().Add(time.Duration(2) * time.Hour)

	claims := jws.Claims{}
	claims.SetIssuer("HyperPaaS")
	claims.SetSubject(strconv.Itoa(user.ID))
	claims.SetExpiration(expiration)
	claims.SetIssuedAt(time.Now().UTC())
	claims.SetJWTID(uuid.NewV4().String())

	token, err := s.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	return &entity.Token{
		Value:   token,
		Expires: expiration.Unix(),
	}, nil
}

// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go docker-compose.yml

func main() {
	_ = container.Get(ServiceLoggerKey).(zerolog.Logger)
	// cfg := container.Get(ServiceConfigKey).(*config.Configuration)

	log.Info().Msg("HyperPaaS Installer")

	docker := container.Get(ServiceDockerKey).(*client.Client)

	info, err := docker.Info(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Get Docker Info")

		return
	}

	if info.Swarm.LocalNodeState != "active" {
		log.Error().Msgf("Swarm Status: %s", info.Swarm.LocalNodeState)

		return
	}

	log.Info().Msgf("Swarm Status: %s", info.Swarm.LocalNodeState)

	if !info.Swarm.ControlAvailable {
		log.Error().Msgf("This node is not a Docker Swarm Manager")

		return
	}

	_, err = loadComposeFile("docker-compose.yml")
	if err != nil {
		log.Error().Err(err).Msg("Load docker-compose.yml")

		return
	}
}

func loadComposeFile(filename string) (*composetypes.Config, error) {
	compose, err := Asset("docker-compose.yml")
	if err != nil {
		return nil, errors.Wrap(err, "Asset docker-compose.yml")
	}

	file, err := loader.ParseYAML(compose)
	if err != nil {
		return nil, errors.Wrap(err, "loader.ParseYAML")
	}

	configFile := &composetypes.ConfigFile{
		Filename: "docker-compose.yml",
		Config:   file,
	}

	var details composetypes.ConfigDetails
	details.WorkingDir = "./"
	details.ConfigFiles = []composetypes.ConfigFile{*configFile}
	details.Environment, err = buildEnvironment(os.Environ())

	config, err := loader.Load(details)
	if err != nil {
		if fpe, ok := err.(*loader.ForbiddenPropertiesError); ok {
			return nil, errors.Wrapf(
				err,
				"Compose file contains unsupported options:\n\n%s\n",
				propertyWarnings(fpe.Properties),
			)
		}

		return nil, errors.Wrap(err, "loader.Load")
	}

	unsupportedProperties := loader.GetUnsupportedProperties(details)
	if len(unsupportedProperties) > 0 {
		return nil, fmt.Errorf(
			"Ignoring unsupported options: %s\n\n",
			strings.Join(unsupportedProperties, ", "),
		)
	}

	deprecatedProperties := loader.GetDeprecatedProperties(details)
	if len(deprecatedProperties) > 0 {
		return nil, fmt.Errorf(
			"Ignoring deprecated options:\n\n%s\n\n",
			propertyWarnings(deprecatedProperties),
		)
	}

	return config, nil
}

func buildEnvironment(env []string) (map[string]string, error) {
	result := make(map[string]string, len(env))

	for _, s := range env {
		// if value is empty, s is like "K=", not "K".
		if !strings.Contains(s, "=") {
			return result, errors.Errorf("unexpected environment %q", s)
		}

		kv := strings.SplitN(s, "=", 2)

		result[kv[0]] = kv[1]
	}

	return result, nil
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string

	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}

	sort.Strings(msgs)

	return strings.Join(msgs, "\n\n")
}

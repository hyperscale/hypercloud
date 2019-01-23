// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package compose

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/docker/cli/cli/compose/loader"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/pkg/errors"
)

// Loader compose file
func Loader(compose []byte, envs map[string]string) (*composetypes.Config, error) {
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

	for key, val := range envs {
		details.Environment[key] = val
	}

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

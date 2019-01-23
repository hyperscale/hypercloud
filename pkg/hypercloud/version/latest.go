// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package version

import (
	"encoding/json"
	"net/http"

	"github.com/coreos/go-semver/semver"
	"github.com/rs/zerolog/log"
)

var githubTagsAPI = "https://api.github.com/repos/hyperscale/hypercloud/tags"

type githubTag struct {
	Name string `json:"name"`
}

// GetLatest version
func GetLatest() (version *semver.Version, err error) {
	resp, err := http.Get(githubTagsAPI)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	tags := []*githubTag{}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	latest := Version

	for _, tag := range tags {
		if tag.Name[0] == 'v' {
			tag.Name = tag.Name[1:]
		}

		tagVersion, err := semver.NewVersion(tag.Name)
		if err != nil {
			log.Error().Err(err).Msgf("Parsing version: %s", tag.Name)

			continue
		}

		if latest.LessThan(*tagVersion) {
			latest = tagVersion
		}
	}

	return latest, nil
}

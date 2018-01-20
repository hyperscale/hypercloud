// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package version

import (
	"encoding/json"
	"net/http"

	"github.com/coreos/go-semver/semver"
)

const githubTagsAPI = "https://api.github.com/repos/hyperscale/hyperpaas/tags"

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
		tagVersion, err := semver.NewVersion(tag.Name)
		if err != nil {
			continue
		}

		if latest.LessThan(*tagVersion) {
			latest = tagVersion
		}
	}

	return latest, nil
}

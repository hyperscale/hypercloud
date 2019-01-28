// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coreos/go-semver/semver"
	"github.com/stretchr/testify/assert"
)

func TestGetLatest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []*githubTag{
			&githubTag{
				Name: "v1.0.0",
			},
			&githubTag{
				Name: "v0.2.0",
			},
			&githubTag{
				Name: "v0.1.0",
			},
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer ts.Close()

	githubTagsAPI = ts.URL

	Version = semver.New("0.2.0")

	v, err := GetLatest()
	assert.NoError(t, err)

	assert.Equal(t, "1.0.0", v.String())
}

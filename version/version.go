// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package version

import "github.com/coreos/go-semver/semver"

var (
	// Tag full version
	Tag = "0.0.0"

	// Revision hash
	Revision = "dev"

	// BuildAt date
	BuildAt string
)

// Version is the specification version that the package types support.
var Version = semver.New(Tag)

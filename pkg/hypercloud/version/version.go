// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package version

import (
	"fmt"
	"runtime"

	"github.com/coreos/go-semver/semver"
)

var version string
var gitCommit string
var gitTreeState string
var buildDate string
var platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

// Info version
type Info struct {
	Version      string
	GitVersion   string
	GitCommit    string
	GitTreeState string
	BuildDate    string
	GoVersion    string
	Compiler     string
	Platform     string
}

// Get returns the version and buildtime information about the binary
func Get() *Info {
	// These variables typically come from -ldflags settings to `go build`
	return &Info{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     platform,
	}
}

// Version is the specification version that the package types support.
var Version = semver.New("0.0.0")

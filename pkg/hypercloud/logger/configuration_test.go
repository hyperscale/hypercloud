// Copyright 2019 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationLevel(t *testing.T) {
	expects := []struct {
		Value    string
		Expected zerolog.Level
	}{
		{
			Value:    "debug",
			Expected: zerolog.DebugLevel,
		},
		{
			Value:    "info",
			Expected: zerolog.InfoLevel,
		},
		{
			Value:    "warn",
			Expected: zerolog.WarnLevel,
		},
		{
			Value:    "error",
			Expected: zerolog.ErrorLevel,
		},
		{
			Value:    "fatal",
			Expected: zerolog.FatalLevel,
		},
		{
			Value:    "panic",
			Expected: zerolog.PanicLevel,
		},
		{
			Value:    "",
			Expected: zerolog.Disabled,
		},
	}

	for _, item := range expects {
		c := &Configuration{
			LevelName: item.Value,
		}

		assert.Equal(t, item.Expected, c.Level())
	}
}

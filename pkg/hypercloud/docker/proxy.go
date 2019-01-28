// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docker

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/docker/go-connections/sockets"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// NewProxy Docker
func NewProxy(host string, base string) (*httputil.ReverseProxy, error) {
	sock, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	target, err := url.Parse("http://docker")
	if err != nil {
		return nil, err
	}

	targetQuery := target.RawQuery

	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(
			target.Path,
			strings.Replace(req.URL.Path, base, "/", 1),
		)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		req.Header.Set("Origin", "http://docker")

		for _, key := range []string{"Access-Control-Request-Headers", "Access-Control-Request-Method"} {
			if val := req.Header.Get(key); val != "" {
				req.Header.Set(key, val)
			}
		}
	}

	transport := new(http.Transport)
	if err := sockets.ConfigureTransport(transport, sock.Scheme, sock.Path); err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}, nil
}

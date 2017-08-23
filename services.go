// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package hyperpaas

import (
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/euskadi31/go-service"
	"github.com/hyperscale/hyperpaas/services"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// Service Container
var Service = service.New()

func init() {
	Service.Set("log", func(c *service.Container) interface{} {
		log := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", "hyperpaas").
			//Str("host", host).
			Logger()

		stdlog.SetFlags(0)
		stdlog.SetOutput(log)

		return log
	})

	Service.Set("middleware", func(c *service.Container) interface{} {
		log := c.Get("log").(zerolog.Logger)

		middleware := alice.New()

		middleware = middleware.Append(hlog.NewHandler(log))
		// Install some provided extra handler to set some request's context fields.
		// Thanks to those handler, all our logs will come with some pre-populated fields.
		middleware = middleware.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}))
		middleware = middleware.Append(hlog.RemoteAddrHandler("ip"))
		middleware = middleware.Append(hlog.UserAgentHandler("user_agent"))
		middleware = middleware.Append(hlog.RefererHandler("referer"))
		middleware = middleware.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

		return middleware
	})

	Service.Set("app", func(c *service.Container) interface{} {
		log := c.Get("log").(zerolog.Logger)

		return services.NewAppService(log)
	})
}

// Run Application
func Run() {
	log := Service.Get("log").(zerolog.Logger)

	app := Service.Get("app").(*services.AppService)

	if err := app.Run(); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

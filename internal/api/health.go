package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/nkbhasker/go-auth-starter/internal/core"
	"github.com/nkbhasker/go-auth-starter/internal/health"
)

type healthHandler struct {
	app core.App
}

func NewHealthHandler(app core.App) *healthHandler {
	return &healthHandler{app: app}
}

func (hh *healthHandler) LiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, hh.app.Check())
	}
}

func (hh *healthHandler) ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := health.NewHealth()
		h.SetStatus(health.HealthStatusUp)
		services := map[string]health.Checker{
			"app":      hh.app,
			"postgres": hh.app.DBStore(),
			"redis":    hh.app.CacheStore(),
		}
		type result struct {
			service string
			health  *health.Health
		}
		ch := make(chan result)
		for k, s := range services {
			go func(k string, s health.Checker) {
				ch <- result{service: k, health: s.Check()}
			}(k, s)
		}

		for range services {
			result := <-ch
			if result.health.Status() != health.HealthStatusUp {
				render.Status(r, http.StatusServiceUnavailable)
				h.SetStatus(result.health.Status())
			}
			h.SetInfo(result.service, result.health.Info())
		}
		close(ch)

		render.JSON(w, r, h)
	}
}

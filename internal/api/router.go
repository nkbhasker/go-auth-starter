package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nkbhasker/go-auth-starter/internal/core"
	"github.com/nkbhasker/go-auth-starter/internal/middleware"
	"github.com/nkbhasker/go-auth-starter/internal/misc"
)

type RouterOptions struct {
	App                    core.App
	JwtHelper              misc.JwtHelper
	OtpGenerateRateLimiter core.RateLimiter
	OtpVerifyRateLimiter   core.RateLimiter
}

func SetupRouter(options RouterOptions) http.Handler {
	healthHandler := NewHealthHandler(options.App)
	authHandler := NewAuthHandler(options.App, options.OtpGenerateRateLimiter, options.OtpVerifyRateLimiter)
	userHandler := NewUserHandler(options.App)
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Group(func(r chi.Router) {
		r.Get("/live", healthHandler.LiveHandler())
		r.Get("/ready", healthHandler.ReadyHandler())
		r.Post("/auth/otp", authHandler.OtpHandler())
		r.Post("/auth/signin", authHandler.SignInHandler())
	})

	router.Group(func(r chi.Router) {
		authInterceptor := middleware.NewAuthInterceptor(options.JwtHelper)
		r.Use(authInterceptor.HandlerFunc)
		r.Get("/user/me", userHandler.MeHandler())
		r.Patch("/user/me", userHandler.UpdateUserHandler())
		r.Put("/user/me/email", userHandler.UpdateEmailHandler())
	})

	return router
}

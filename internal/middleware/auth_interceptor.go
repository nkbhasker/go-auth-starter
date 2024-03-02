package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/nkbhasker/go-auth-starter/internal/core"
	"github.com/nkbhasker/go-auth-starter/internal/misc"
)

type authInterceptor struct {
	jwtHelper misc.JwtHelper
}

const (
	Bearer string = "bearer"
)

func NewAuthInterceptor(jwtHelper misc.JwtHelper) *authInterceptor {
	return &authInterceptor{
		jwtHelper: jwtHelper,
	}
}

func (a *authInterceptor) HandlerFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		identity, err := func() (core.Identity, error) {
			accessToken, err := extractTokenFromHeader(r.Header)
			if err != nil {
				return nil, err
			}
			claims, err := a.jwtHelper.VerifyAccessToken(accessToken)
			if err != nil {
				return nil, err
			}
			return core.NewIdentity(claims.ID, claims.Subject)
		}()
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(core.IdentityToContext(ctx, identity)))
	})
}

func extractTokenFromHeader(h http.Header) (string, error) {
	authHeader := h.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != Bearer {
		return "", errors.New("authorization header format must be bearer {token}")
	}

	return authHeaderParts[1], nil
}

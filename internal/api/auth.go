package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/muhlemmer/httpforwarded"

	"github.com/go-chi/render"
	"github.com/nkbhasker/go-auth-starter/internal/core"
	"github.com/nkbhasker/go-auth-starter/internal/misc"
	"github.com/nkbhasker/go-auth-starter/internal/model"
	"github.com/nkbhasker/go-auth-starter/internal/repo"
)

type OtpScopeEnum string

const (
	OtpScopeSignIn      OtpScopeEnum = "SIGN_IN"
	OtpScopeEmailUpdate OtpScopeEnum = "EMAIL_UPDATE"
)

type authHandler struct {
	app                    core.App
	otpGenerateRateLimiter core.RateLimiter
	otpVerifyRateLimiter   core.RateLimiter
}

type otpRequestBody struct {
	Email string `json:"email" validate:"required,email"`
	Scope string `json:"scope" validate:"oneof=SIGN_IN EMAIL_UPDATE"`
}

type signInRequestBody struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

func NewAuthHandler(app core.App, otpGenerateRateLimiter core.RateLimiter, otpVerifyRateLimiter core.RateLimiter) *authHandler {
	return &authHandler{
		app:                    app,
		otpGenerateRateLimiter: otpGenerateRateLimiter,
		otpVerifyRateLimiter:   otpVerifyRateLimiter,
	}
}

func (h *authHandler) OtpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() error {
			ok, err := h.otpGenerateRateLimiter.Evaluate(GetIP(r))
			if !ok || err != nil {
				return fmt.Errorf("too many otp requests")
			}
			otpBody := &otpRequestBody{}
			err = json.NewDecoder(r.Body).Decode(otpBody)
			if err != nil {
				return err
			}
			if otpBody.Scope == "" {
				otpBody.Scope = string(OtpScopeSignIn)
			}
			err = h.app.Validate().Struct(otpBody)
			if err != nil {
				return err
			}
			otp, err := misc.GenerateOtp()
			if err != nil {
				return err
			}
			err = h.app.Emailer().SendSignInOTP(otpBody.Email, otp)
			if err != nil {
				return err
			}
			key := fmt.Sprintf("%s_%s_%s", otpBody.Scope, repo.AuthKeyOTP, otpBody.Email)
			err = h.app.Repo().AuthRepo().SaveOTP(r.Context(), key, otp)
			if err != nil {
				return err
			}
			// Reset otp verify rate limit
			return h.otpVerifyRateLimiter.Reset(otpBody.Email)
		}()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		render.JSON(w, r, map[string]interface{}{
			"success": true,
		})
	}
}

func (h *authHandler) SignInHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := func() (string, error) {
			signInBody := &signInRequestBody{}
			err := json.NewDecoder(r.Body).Decode(signInBody)
			if err != nil {
				return "", err
			}
			err = h.app.Validate().Struct(signInBody)
			if err != nil {
				return "", err
			}
			ok, err := h.otpVerifyRateLimiter.Evaluate(signInBody.Email)
			if !ok || err != nil {
				return "", fmt.Errorf("too many invalid otp attempts")
			}
			key := fmt.Sprintf("%s_%s_%s", OtpScopeSignIn, repo.AuthKeyOTP, signInBody.Email)
			otp, err := h.app.Repo().AuthRepo().GetOTP(r.Context(), key)
			if err != nil {
				return "", err
			}
			if ok := misc.ValidateOtp(otp, signInBody.OTP); !ok {
				return "", fmt.Errorf("invalid otp")
			}
			user, err := h.app.Repo().UserRepo().GetByEmail(signInBody.Email)
			// Create new user
			if errors.Is(err, repo.ErrUserNotFound) {
				user, err = h.app.Repo().UserRepo().New(model.User{Email: &signInBody.Email, IsEmailVerified: true})
				if err != nil {
					return "", err
				}
				err = h.app.Repo().UserRepo().Create(user)
				if err != nil {
					return "", err
				}
			}
			if err != nil {
				return "", err
			}

			return h.app.Repo().AccessTokenRepo().Create(user.ID.String(), signInBody.Email)
		}()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"success":     true,
			"accessToken": accessToken,
		})
	}
}

func GetIP(r *http.Request) string {
	fwd, err := httpforwarded.ParseFromRequest(r)
	if err == nil && len(fwd["X-Forwarded-For"]) != 0 {
		return fwd["X-Forwarded-For"][0]
	}

	return r.RemoteAddr
}

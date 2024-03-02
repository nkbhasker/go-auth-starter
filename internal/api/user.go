package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/nkbhasker/go-auth-starter/internal/core"
	"github.com/nkbhasker/go-auth-starter/internal/misc"
	"github.com/nkbhasker/go-auth-starter/internal/model"
	"github.com/nkbhasker/go-auth-starter/internal/repo"
)

type userHandler struct {
	app core.App
}

type updateUserRequestBody struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender" validate:"oneof=MALE FEMALE"`
}

type updateEmailRequestBody struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func NewUserHandler(app core.App) *userHandler {
	return &userHandler{app: app}
}

func (h *userHandler) MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity := core.IdentityFromContext(r.Context())
		user, err := h.app.Repo().UserRepo().Get(identity.UserID())
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
			"user":    user,
		})
	}
}

func (h *userHandler) UpdateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := func() (*model.User, error) {
			identity := core.IdentityFromContext(r.Context())
			updateUserBody := &updateUserRequestBody{}
			err := json.NewDecoder(r.Body).Decode(updateUserBody)
			if err != nil {
				return nil, err
			}
			user, err := h.app.Repo().UserRepo().Get(identity.UserID())
			if err != nil {
				return nil, err
			}
			err = h.app.Repo().UserRepo().Update(user)
			if err != nil {
				return nil, err
			}

			return user, nil
		}()
		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		render.JSON(w, r, map[string]interface{}{
			"success": true,
			"user":    user,
		})
	}
}

func (h *userHandler) UpdateEmailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := func() error {
			identity := core.IdentityFromContext(r.Context())
			updateEmailBody := &updateEmailRequestBody{}
			err := json.NewDecoder(r.Body).Decode(updateEmailBody)
			if err != nil {
				return err
			}
			key := fmt.Sprintf("%s_%s_%s", OtpScopeEmailUpdate, repo.AuthKeyOTP, updateEmailBody.Email)
			otp, err := h.app.Repo().AuthRepo().GetOTP(r.Context(), key)
			if err != nil {
				return err
			}
			if ok := misc.ValidateOtp(otp, updateEmailBody.OTP); !ok {
				return fmt.Errorf("invalid otp")
			}
			user, err := h.app.Repo().UserRepo().Get(identity.UserID())
			if err != nil {
				return nil
			}
			user.Email = &updateEmailBody.Email
			err = h.app.Repo().UserRepo().Update(user)
			if err != nil {
				return nil
			}

			return nil
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

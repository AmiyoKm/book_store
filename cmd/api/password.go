package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	mailer "github.com/AmiyoKm/book_store/internal/mail"
	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/google/uuid"
)

type passwordResetRequestPayload struct {
	Email string `json:"email" validate:"required,email"`
}
type TokenResponse struct {
	Token string `json:"token"`
}

// passwordResetRequestHandler godoc
//
//	@Summary		Send Reset Password Request
//	@Description	Send a Reset Password Request by sending a mail to the user
//	@Tags			password
//	@Accept			json
//	@Produce		json
//
//	@Param			payload	body		passwordResetRequestPayload	true	"Password Reset Request Payload"
//	@Success		201		{object}	TokenResponse				"Reset Password Request Response"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/password/reset-request [post]
func (app *Application) passwordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	var payload passwordResetRequestPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	ID, err := app.store.Users.CreatePasswordRequest(ctx, user, hashToken, app.cfg.mail.exp)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	isProdEnv := app.cfg.env == "PRODUCTION"

	PasswordResetURL := fmt.Sprintf("%s/reset-password?token=%s", app.cfg.frontendURL, plainToken)

	vars := struct {
		Username         string
		PasswordResetURL string
	}{
		Username:         user.Username,
		PasswordResetURL: PasswordResetURL,
	}

	_, err = app.mail.SendPasswordRequestMail(mailer.PasswordChangeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending password request email", "email", err)
		if err := app.store.Users.DeletePasswordRequest(ctx, *ID); err != nil {
			app.logger.Errorw("error deleting password request", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infof("Sending email from: %s", app.cfg.mail.fromEmail)
	if err := jsonResponse(w, http.StatusCreated, TokenResponse{plainToken}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type PasswordResetVerifyResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// passwordRequestVerifyHandler godoc
//
//	@Summary		Verify the password reset request
//	@Description	Verifies the password reset request sent by the email
//	@Tags			password
//	@Accept			json
//	@Produce		json
//	@Param			token	query		string						true	"Password reset token"
//	@Success		200		{object}	PasswordResetVerifyResponse	"Reset Password Request Response"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/password/request/verify [get]
func (app *Application) passwordRequestVerifyHandler(w http.ResponseWriter, r *http.Request) {
	plainToken := r.URL.Query().Get("token")

	if plainToken == "" {
		app.badRequestError(w, r, fmt.Errorf("token is required"))
		return
	}
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	ctx := r.Context()

	request, err := app.store.Users.GetPasswordRequest(ctx, hashToken)
	if err != nil {
		if errors.Is(err, store.ErrorNotFound) {
			app.notFoundError(w, r, fmt.Errorf("invalid or expired token"))
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}

	if request.Used {
		app.badRequestError(w, r, fmt.Errorf("token already used"))
		return
	}

	if time.Now().After(request.Expiry) {
		app.badRequestError(w, r, fmt.Errorf("token expired"))
		return
	}
	response := PasswordResetVerifyResponse{
		"Token is valid",
		fmt.Sprintf("%d", request.UserID),
	}

	if err := jsonResponse(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
	}
}

type passwordResetPayload struct {
	UserID      int    `json:"user_id" validate:"required"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=5"`
}
type passwordResetResponse struct {
	Message string `json:"message"`
}

// passwordRequestVerifyHandler godoc
//
//	@Summary		Reset the password
//	@Description	Reset the password
//	@Tags			password
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		passwordResetPayload	true	"Reset Password Payload"
//	@Success		200		{object}	passwordResetResponse	"Reset Password Request Response"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/password/reset [post]
func (app *Application) passwordResetHandler(w http.ResponseWriter, r *http.Request) {
	var payload passwordResetPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	hash := sha256.Sum256([]byte(payload.Token))
	hashToken := hex.EncodeToString(hash[:])

	ctx := r.Context()
	request, err := app.store.Users.GetPasswordRequest(ctx, hashToken)
	if err != nil {
		if errors.Is(err, store.ErrorNotFound) {
			app.notFoundError(w, r, fmt.Errorf("invalid or expired token"))
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}
	if request.UserID != payload.UserID {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized access"))
		return
	}

	hashedPassword := &store.Password{}
	hashedPassword.Set(payload.NewPassword)

	if err := app.store.Users.UpdatePassword(ctx, payload.UserID, hashedPassword); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.store.Users.MarkPasswordRequestAsUsed(ctx, hashToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infof("Password successfully reset for user ID: %d", payload.UserID)
	jsonResponse(w, http.StatusOK, passwordResetResponse{"Password updated successfully"})
}

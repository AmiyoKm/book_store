package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	mailer "github.com/AmiyoKm/book_store/internal/mail"
	"github.com/google/uuid"
)

func (app *Application) passwordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	ctx := r.Context()
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
	if err := jsonResponse(w, http.StatusCreated, plainToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

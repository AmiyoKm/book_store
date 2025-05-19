package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	mailer "github.com/AmiyoKm/book_store/internal/mail"
	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type createUserPayload struct {
	Username string `json:"username" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=255,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=user moderator admin"`
}

type UserWithToken struct {
	User  *store.User
	Token string `json:"token"`
}
// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *Application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload createUserPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()

	hashedPass := store.Password{}
	hashedPass.Set(payload.Password)
	if payload.Role == "" {
		payload.Role = "user"
	}
	role := store.Role{Name: payload.Role}
	user := &store.User{
		Username: payload.Username,
		Password: hashedPass,
		Email:    payload.Email,
		Role:     role,
	}
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.cfg.mail.exp)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
			return
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	isProdEnv := app.cfg.env == "PRODUCTION"

	ActivationURL := fmt.Sprintf("%s/confirm/%s", app.cfg.frontendURL, plainToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: ActivationURL,
	}
	_, err = app.mail.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "email", err)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infof("Sending email from: %s", app.cfg.mail.fromEmail)
	if err := jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type loginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		loginUserPayload	true	"User credentials"
//	@Success		200		{string}	string				"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *Application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload loginUserPayload

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
			app.unauthorizedError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := user.Password.ComparePassword(payload.Password); err != nil {
		app.unauthorizedError(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.cfg.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.cfg.auth.token.iss,
		"aud": app.cfg.auth.token.iss,
	}

	token, err := app.auth.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, token); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
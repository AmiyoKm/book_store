package main

import (
	"net/http"

	"github.com/AmiyoKm/book_store/internal/store"
)

type createUserPayload struct {
	Username string `json:"username" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,max=255,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=user moderator admin"`
}

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
	err := app.store.Users.Create(ctx, user)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type loginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

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

	if err := user.Password.ComparePassword(payload.Password) ; err != nil {
		app.unauthorizedError(w,r,err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

package main

import (
	"fmt"
	"net/http"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
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
	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if user == nil {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized error"))
		return
	}
	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type updateUserPayload struct {
	UserName *string `json:"username" validate:"required,min=2"`
}

func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if user == nil {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized error"))
		return
	}

	var payload updateUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if payload.UserName != nil {
		user.Username = *payload.UserName
	}
	err := app.store.Users.Update(r.Context(), user)
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
	if err := jsonResponse(w, http.StatusAccepted, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}

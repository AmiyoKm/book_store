package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

type Review struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    BookID    int       `json:"book_id"`
    Content   string    `json:"content"`
    Rating    int       `json:"rating"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
// activateUserHandler godoc
//	@Summary		Activate user account
//	@Description	Activates a user account using a provided token.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string	true	"Activation Token"
//	@Success		200		{object}	string	"User account activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/activate/{token} [get]
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
	if err := jsonResponse(w, http.StatusOK, map[string]string{"message": "User account activated successfully"}); err != nil {
		app.internalServerError(w, r, err)
	}
}
// getUserHandler godoc
//	@Summary		Get current user details
//	@Description	Retrieves details of the authenticated user.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	store.User	"User details"
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/me [get]
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
// updateUserHandler godoc
//	@Summary		Update user details
//	@Description	Updates the authenticated user's details.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		updateUserPayload	true	"User update payload"
//	@Success		202		{object}	store.User			"User updated"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/me [patch]
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
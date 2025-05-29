package main

import (
	"net/http"
	"strconv"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

type addToWishlistPayload struct {
	BookID int `json:"book_id" validate:"required"`
}

// addToWishlistHandler godoc
//
//	@Summary		Add Book to Wishlist
//	@Description	Add Book to Wishlist
//	@Tags			wishlist
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		addToWishlistPayload	true	"Book ID"
//	@Success		201		{object}	store.Wishlist			"Create Wishlist"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist [post]
func (app *Application) addToWishlistHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	var payload addToWishlistPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	wishlist := &store.Wishlist{
		UserID: user.ID,
		BookID: payload.BookID,
	}
	ctx := r.Context()
	err := app.store.WishLists.Create(ctx, wishlist)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusCreated, wishlist); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getWishlistHandler godoc
//
//	@Summary		Get Book Wishlist
//	@Description	Get Book Wishlist
//	@Tags			wishlist
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Book	"Book List"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist [get]
func (app *Application) getWishlistHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	books, err := app.store.WishLists.GetWishlistBooks(r.Context(), user.ID)

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
	if err := jsonResponse(w, http.StatusOK, books); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type removeWishlistResponse struct {
	BookID  int    `json:"book_id"`
	Message string `json:"message"`
}

// removeFromWishlistHandler godoc
//
//	@Summary		Delete Book From Wishlist
//	@Description	Delete Book From Wishlist
//	@Tags			wishlist
//	@Accept			json
//	@Produce		json
//	@Param			bookID	path		int						true	"BOOK ID"
//	@Success		200		{object}	removeWishlistResponse	"DELETED RESPONSE"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/wishlist/{bookID} [delete]
func (app *Application) removeFromWishlistHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	paramID := chi.URLParam(r, "bookID")
	bookID, err := strconv.Atoi(paramID)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	err = app.store.WishLists.Delete(r.Context(), user.ID, bookID)
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
	res := removeWishlistResponse{
		Message: "Removed from wishlist",
		BookID:  bookID,
	}
	if err := jsonResponse(w, http.StatusOK, res); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
